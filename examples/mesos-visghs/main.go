/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/tcolgate/vonq/ghs"
	"github.com/tcolgate/vonq/graph"
	"github.com/tcolgate/vonq/graphalg"
	"github.com/tcolgate/vonq/tracer"

	"github.com/gogo/protobuf/proto"
	log "github.com/golang/glog"
	"github.com/mesos/mesos-go/auth"
	"github.com/mesos/mesos-go/auth/sasl"
	"github.com/mesos/mesos-go/auth/sasl/mech"
	exec "github.com/mesos/mesos-go/executor"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
	"golang.org/x/net/context"
)

func init() {
	flag.Parse()
}

const (
	CPUS_PER_EXECUTOR   = 0.01
	CPUS_PER_TASK       = 0.1
	MEM_PER_EXECUTOR    = 64
	MEM_PER_TASK        = 64
	PORTS_PER_TASK      = 1
	defaultArtifactPort = 12345
)

var (
	address      = flag.String("address", "127.0.0.1", "Binding address for artifact server")
	artifactPort = flag.Int("artifactPort", defaultArtifactPort, "Binding port for artifact server")
	authProvider = flag.String("mesos_authentication_provider", sasl.ProviderName,
		fmt.Sprintf("Authentication provider to use, default is SASL that supports mechanisms: %+v", mech.ListSupported()))
	master     = flag.String("master", "127.0.0.1:5050", "Master address <ip:port>")
	ghsNode    = flag.Bool("node", false, "Run a ghs node")
	nodeCount  = flag.Int("node-count", 5, "Total task count to run.")
	tracerAddr = flag.String("tracerAddr", "", "Tracer GRPC url")

	mesosAuthPrincipal  = flag.String("mesos_authentication_principal", "", "Mesos authentication principal.")
	mesosAuthSecretFile = flag.String("mesos_authentication_secret_file", "", "Mesos authentication secret file.")
)

type visghsScheduler struct {
	tracer.Tracer

	nexec *mesos.ExecutorInfo

	nodesErrored  int
	nodesLaunched int
	nodesFinished int
	nodeTasks     int
}

func newvisghsScheduler(nexec *mesos.ExecutorInfo) *visghsScheduler {
	total := *nodeCount
	return &visghsScheduler{
		nexec:     nexec,
		nodeTasks: total,
	}
}

func (s *visghsScheduler) Registered(driver sched.SchedulerDriver, frameworkId *mesos.FrameworkID, masterInfo *mesos.MasterInfo) {
	log.Infoln("Framework Registered with Master ", masterInfo)
}

func (s *visghsScheduler) Reregistered(driver sched.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	log.Infoln("Framework Re-Registered with Master ", masterInfo)
	_, err := driver.ReconcileTasks([]*mesos.TaskStatus{})
	if err != nil {
		log.Errorf("failed to request task reconciliation: %v", err)
	}
}

func (s *visghsScheduler) Disconnected(sched.SchedulerDriver) {
	log.Warningf("disconnected from master")
}

func (s *visghsScheduler) ResourceOffers(driver sched.SchedulerDriver, offers []*mesos.Offer) {
	if (s.nodesLaunched - s.nodesErrored) >= s.nodeTasks {
		log.Info("decline all of the offers since all of our tasks are already launched")
		ids := make([]*mesos.OfferID, len(offers))
		for i, offer := range offers {
			ids[i] = offer.Id
		}
		driver.LaunchTasks(ids, []*mesos.TaskInfo{}, &mesos.Filters{RefuseSeconds: proto.Float64(120)})
		return
	}
	for _, offer := range offers {
		cpuResources := util.FilterResources(offer.Resources, func(res *mesos.Resource) bool {
			return res.GetName() == "cpus"
		})
		cpus := 0.0
		for _, res := range cpuResources {
			cpus += res.GetScalar().GetValue()
		}

		memResources := util.FilterResources(offer.Resources, func(res *mesos.Resource) bool {
			return res.GetName() == "mem"
		})
		mems := 0.0
		for _, res := range memResources {
			mems += res.GetScalar().GetValue()
		}

		portResources := util.FilterResources(offer.Resources, func(res *mesos.Resource) bool {
			return res.GetName() == "ports"
		})
		ports := []*mesos.Value_Range{}
		portsCount := uint64(0)
		for _, res := range portResources {
			for _, rs := range res.GetRanges().GetRange() {
				ports = append(ports, rs)
				portsCount += 1 + rs.GetEnd() - rs.GetBegin()
			}
		}

		log.Infoln("Received Offer <", offer.Id.GetValue(), "> with cpus=", cpus, " mem=", mems, " ports=", ports)

		remainingCpus := cpus
		remainingMems := mems
		remainingPorts := ports
		remainingPortsCount := portsCount

		// account for executor resources if there's not an executor already running on the slave
		if len(offer.ExecutorIds) == 0 {
			remainingCpus -= CPUS_PER_EXECUTOR
			remainingMems -= MEM_PER_EXECUTOR
		}

		var tasks []*mesos.TaskInfo
		for (s.nodesLaunched-s.nodesErrored) < s.nodeTasks &&
			CPUS_PER_TASK <= remainingCpus &&
			MEM_PER_TASK <= remainingMems &&
			PORTS_PER_TASK <= remainingPortsCount {
			log.Infoln("Ports <", remainingPortsCount, remainingPorts)

			s.nodesLaunched++

			taskId := &mesos.TaskID{
				Value: proto.String(strconv.Itoa(s.nodesLaunched)),
			}

			taskPorts := []*mesos.Value_Range{}
			leftOverPorts := []*mesos.Value_Range{}
			for t := 0; t < PORTS_PER_TASK; t++ {
				if len(remainingPorts) < 1 {
					// failed to allocate port, oh no!
				}
				ps := remainingPorts[0]
				//take the first port from the first
				pb := ps.GetBegin()
				pe := ps.GetEnd()
				//Create one range per port we need, it's easier this way
				tp := mesos.Value_Range{}
				p := pb
				tp.Begin = &p
				tp.End = &p
				taskPorts = append(taskPorts, &tp)

				pb++
				if pb <= pe {
					rpb := pb
					rpe := pe
					rtp := mesos.Value_Range{Begin: &rpb, End: &rpe}
					leftOverPorts = append(leftOverPorts, &rtp)
				}
				for _, ps := range remainingPorts[1:] {
					leftOverPorts = append(leftOverPorts, ps)
				}
			}

			vonqPort := (uint32)(taskPorts[0].GetBegin())
			task := &mesos.TaskInfo{
				Name:     proto.String("visghs-node-" + taskId.GetValue()),
				TaskId:   taskId,
				SlaveId:  offer.SlaveId,
				Executor: s.nexec,
				Discovery: &mesos.DiscoveryInfo{
					Name: proto.String("visghs"),
					//Visibility: mesos.DiscoveryInfo_EXTERNAL.Enum(),
					Visibility: mesos.DiscoveryInfo_FRAMEWORK.Enum(),
					Ports: &mesos.Ports{
						Ports: []*mesos.Port{
							{Protocol: proto.String("UDP"),
								Visibility: mesos.DiscoveryInfo_EXTERNAL.Enum(),
								Name:       proto.String("udpprobe"),
								Number:     proto.Uint32(vonqPort)},
							{Protocol: proto.String("TCP"),
								Visibility: mesos.DiscoveryInfo_EXTERNAL.Enum(),
								Name:       proto.String("vonqrpc"),
								Number:     proto.Uint32(vonqPort)},
						},
					},
				},
				Resources: []*mesos.Resource{
					util.NewScalarResource("cpus", CPUS_PER_TASK),
					util.NewScalarResource("mem", MEM_PER_TASK),
					util.NewRangesResource("ports", taskPorts),
				},
			}
			log.Infof("Prepared task: %s with offer %s for launch\n", task.GetName(), offer.Id.GetValue())

			tasks = append(tasks, task)
			remainingCpus -= CPUS_PER_TASK
			remainingMems -= MEM_PER_TASK
			remainingPorts = leftOverPorts
			remainingPortsCount -= PORTS_PER_TASK
		}
		log.Infoln("Launching ", len(tasks), "tasks for offer", offer.Id.GetValue())
		driver.LaunchTasks([]*mesos.OfferID{offer.Id}, tasks, &mesos.Filters{RefuseSeconds: proto.Float64(5)})
	}
}

func (sched *visghsScheduler) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {
	log.Infoln("Status update: task", status.TaskId.GetValue(), " is in state ", status.State.Enum().String())
	if status.GetState() == mesos.TaskState_TASK_FINISHED {
		sched.nodesFinished++
		driver.ReviveOffers() // TODO(jdef) rate-limit this
	}

	if sched.nodesFinished >= sched.nodeTasks {
		log.Infoln("Total tasks completed, stopping framework.")
		driver.Stop(false)
	}

	if status.GetState() == mesos.TaskState_TASK_LOST ||
		status.GetState() == mesos.TaskState_TASK_KILLED ||
		status.GetState() == mesos.TaskState_TASK_FAILED ||
		status.GetState() == mesos.TaskState_TASK_ERROR {
		sched.nodesErrored++
	}
}

func (s *visghsScheduler) OfferRescinded(_ sched.SchedulerDriver, oid *mesos.OfferID) {
	log.Errorf("offer rescinded: %v", oid)
}
func (s *visghsScheduler) FrameworkMessage(_ sched.SchedulerDriver, eid *mesos.ExecutorID, sid *mesos.SlaveID, msg string) {
	log.Errorf("framework message from executor %q slave %q: %q", eid, sid, msg)
}
func (s *visghsScheduler) SlaveLost(_ sched.SchedulerDriver, sid *mesos.SlaveID) {
	log.Errorf("slave lost: %v", sid)
}
func (s *visghsScheduler) ExecutorLost(_ sched.SchedulerDriver, eid *mesos.ExecutorID, sid *mesos.SlaveID, code int) {
	log.Errorf("executor %q lost on slave %q code %d", eid, sid, code)
}
func (s *visghsScheduler) Error(_ sched.SchedulerDriver, err string) {
	log.Errorf("Scheduler received error: %v", err)
}

func (s *visghsScheduler) OnRun() {
}

// returns (downloadURI, basename(path))
func serveSelf() *string {
	http.HandleFunc("/executor", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/proc/self/exe")
	})

	hostURI := fmt.Sprintf("http://%s:%d/%s", *address, *artifactPort, "/executor")
	log.V(2).Infof("Hosting artifact '%s'", hostURI)

	return &hostURI
}

func prepareExecutorInfo(gt net.Addr) *mesos.ExecutorInfo {
	executorUris := []*mesos.CommandInfo_URI{}
	uri := serveSelf()
	executorUris = append(executorUris, &mesos.CommandInfo_URI{Value: uri, Executable: proto.Bool(true)})

	// forward the value of the scheduler's -v flag to the executor
	v := 0
	if f := flag.Lookup("v"); f != nil && f.Value != nil {
		if vstr := f.Value.String(); vstr != "" {
			if vi, err := strconv.ParseInt(vstr, 10, 32); err == nil {
				v = int(vi)
			}
		}
	}
	nodeCommand := fmt.Sprintf("./executor -logtostderr=true -v=%d -node -tracerAddr %s", v, gt.String())
	log.V(2).Info("nodeCommand: ", nodeCommand)

	// Create mesos scheduler driver.
	return &mesos.ExecutorInfo{
		ExecutorId: util.NewExecutorID("default"),
		Name:       proto.String("visghs-node"),
		Source:     proto.String("visghs"),
		Command: &mesos.CommandInfo{
			Value: proto.String(nodeCommand),
			Uris:  executorUris,
		},
		Resources: []*mesos.Resource{
			util.NewScalarResource("cpus", CPUS_PER_EXECUTOR),
			util.NewScalarResource("mem", MEM_PER_EXECUTOR),
		},
	}
}

func parseIP(address string) net.IP {
	addr, err := net.LookupIP(address)
	if err != nil {
		log.Fatal(err)
	}
	if len(addr) < 1 {
		log.Fatalf("failed to parse IP from address '%v'", address)
	}
	return addr[0]
}

type ghsvisExecutor struct {
	tasksLaunched int
	*tracer.Tracer
	*grpc.Server
}

func newVISGHSExecutor() *ghsvisExecutor {
	log.Infoln("Initializing the VISGHS Executor...", *tracerAddr)
	c, err := tracer.NewGRPCDisplayClient(*tracerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	t := tracer.New(c)

	// Establish gRPC on the deisgnated port
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("grpc failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	// build command executor
	nexec := prepareExecutorInfo(lis.Addr())
	s := newvisghsScheduler(nexec)

	t := tracer.NewHTTPDisplay(http.DefaultServeMux, s.OnRun)
	server := tracer.NewGRPCServer(t)
	tracer.RegisterTraceServiceServer(grpcServer, server)
	go grpcServer.Serve(lis)
	time.Sleep(2 * time.Second)
	go http.ListenAndServe(":12345", nil)

	return &ghsvisExecutor{tasksLaunched: 0, Tracer: t}
}

func (exec *ghsvisExecutor) Registered(driver exec.ExecutorDriver, execInfo *mesos.ExecutorInfo, fwinfo *mesos.FrameworkInfo, slaveInfo *mesos.SlaveInfo) {
	fmt.Println("Registered Executor on slave ", slaveInfo.GetHostname())
}

func (exec *ghsvisExecutor) Reregistered(driver exec.ExecutorDriver, slaveInfo *mesos.SlaveInfo) {
	fmt.Println("Re-registered Executor on slave ", slaveInfo.GetHostname())
}

func (exec *ghsvisExecutor) Disconnected(exec.ExecutorDriver) {
	fmt.Println("Executor disconnected.")
}

func (exec *ghsvisExecutor) LaunchTask(driver exec.ExecutorDriver, taskInfo *mesos.TaskInfo) {
	fmt.Println("Launching task", taskInfo.GetName(), "with command", taskInfo.Command.GetValue())
	fmt.Println(taskInfo)

	runStatus := &mesos.TaskStatus{
		TaskId: taskInfo.GetTaskId(),
		State:  mesos.TaskState_TASK_RUNNING.Enum(),
	}
	_, err := driver.SendStatusUpdate(runStatus)
	if err != nil {
		fmt.Println("Got error", err)
	}

	exec.tasksLaunched++
	fmt.Println("Total tasks launched ", exec.tasksLaunched)
	//
	// this is where one would perform the requested task
	//
	finishTask := func() {
		// finish task
		fmt.Println("Finishing task", taskInfo.GetName())
		finStatus := &mesos.TaskStatus{
			TaskId: taskInfo.GetTaskId(),
			State:  mesos.TaskState_TASK_FINISHED.Enum(),
		}
		if _, err := driver.SendStatusUpdate(finStatus); err != nil {
			fmt.Println("error sending FINISHED", err)
		}
		fmt.Println("Task finished", taskInfo.GetName())
	}

	starting := &mesos.TaskStatus{
		TaskId: taskInfo.GetTaskId(),
		State:  mesos.TaskState_TASK_STARTING.Enum(),
	}
	if _, err := driver.SendStatusUpdate(starting); err != nil {
		fmt.Println("error sending STARTING", err)
	}
	go func() {
		running := &mesos.TaskStatus{
			TaskId: taskInfo.GetTaskId(),
			State:  mesos.TaskState_TASK_RUNNING.Enum(),
		}
		if _, err := driver.SendStatusUpdate(running); err != nil {
			fmt.Println("error sending RUNNING", err)
		}
		for {
			select {
			case <-time.Tick(5 * time.Second):
				fmt.Println("In loop", taskInfo.String)
				exec.Log(graph.GraphID{}, graph.AlgorithmID{}, taskInfo.GetTaskId().String(), "Hello")
			}
		}
		finishTask()
	}()
}

func (exec *ghsvisExecutor) KillTask(exec.ExecutorDriver, *mesos.TaskID) {
	fmt.Println("Kill task")
}

func (exec *ghsvisExecutor) FrameworkMessage(driver exec.ExecutorDriver, msg string) {
	fmt.Println("Got framework message: ", msg)
}

func (exec *ghsvisExecutor) Shutdown(exec.ExecutorDriver) {
	fmt.Println("Shutting down the executor")
}

func (exec *ghsvisExecutor) Error(driver exec.ExecutorDriver, err string) {
	fmt.Println("Got error message:", err)
}

type state struct {
	ghs []*ghs.State
	wg  *sync.WaitGroup
}

func (s *state) OnRun() {
	s.ghs[0].WakeUp()
	s.wg.Wait()

	log.Infoln("finished")
}

func setupGHS(mux *http.ServeMux) {
	s := state{}

	t := tracer.New(tracer.NewHTTPDisplay(mux, s.OnRun))

	s.wg = &sync.WaitGroup{}

	count := 25
	s.wg.Add(count)
	nodes := make([]*graphalg.Node, 0)
	s.ghs = make([]*ghs.State, 0)

	for i := 0; i < count; i++ {
		n := &graphalg.Node{
			ID:     graphalg.NodeID(fmt.Sprintf("n%v", i+1)),
			Tracer: t,
		}
		nodes = append(nodes, n)
		s.ghs = append(s.ghs, &ghs.State{Node: n})
	}

	for i := 0; i < count; i++ {
		for j := i + 1; j < count; j++ {
			log.Infoln("join", i, j)
			w := rand.NormFloat64()*1.0 + 4.0
			nodes[i].Join(nodes[j], w, graphalg.MakeChanPair)
		}
	}

	for _, g := range s.ghs {
		go func(g *ghs.State) {
			n := g.Node
			n.NodeUpdate("from node")
			for _, e := range n.Edges() {
				e.EdgeUpdate()
			}

			graphalg.Run(g, s.wg.Done)
		}(g)
	}

	log.Infoln("running")
}

// ----------------------- func main() ------------------------- //

func main() {
	if *ghsNode {
		ghsNodeMain()
		return
	}

	// Allocate a port for GRPC
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("grpc failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	// build command executor
	nexec := prepareExecutorInfo(lis.Addr())
	s := newvisghsScheduler(nexec)

	t := tracer.NewHTTPDisplay(http.DefaultServeMux, s.OnRun)
	server := tracer.NewGRPCServer(t)
	tracer.RegisterTraceServiceServer(grpcServer, server)
	go grpcServer.Serve(lis)
	time.Sleep(2 * time.Second)
	go http.ListenAndServe(":12345", nil)

	// the framework
	fwinfo := &mesos.FrameworkInfo{
		User: proto.String(""), // Mesos-go will fill in user.
		Name: proto.String("visghs"),
	}

	cred := (*mesos.Credential)(nil)
	if *mesosAuthPrincipal != "" {
		fwinfo.Principal = proto.String(*mesosAuthPrincipal)
		cred = &mesos.Credential{
			Principal: proto.String(*mesosAuthPrincipal),
		}
		if *mesosAuthSecretFile != "" {
			_, err := os.Stat(*mesosAuthSecretFile)
			if err != nil {
				log.Fatal("missing secret file: ", err.Error())
			}
			secret, err := ioutil.ReadFile(*mesosAuthSecretFile)
			if err != nil {
				log.Fatal("failed to read secret file: ", err.Error())
			}
			cred.Secret = proto.String(string(secret))
		}
	}

	bindingAddress := parseIP(*address)
	config := sched.DriverConfig{
		Scheduler:      s,
		Framework:      fwinfo,
		Master:         *master,
		Credential:     cred,
		BindingAddress: bindingAddress,
		WithAuthContext: func(ctx context.Context) context.Context {
			ctx = auth.WithLoginProvider(ctx, *authProvider)
			ctx = sasl.WithBindingAddress(ctx, bindingAddress)
			return ctx
		},
	}
	driver, err := sched.NewMesosSchedulerDriver(config)

	if err != nil {
		log.Errorln("Unable to create a SchedulerDriver ", err.Error())
	}

	if stat, err := driver.Run(); err != nil {
		log.Infof("Framework stopped with status %s and error: %s\n", stat.String(), err.Error())
		time.Sleep(2 * time.Second)
		os.Exit(1)
	}
	log.Infof("framework terminating")
}

func ghsNodeMain() {
	fmt.Println("Starting GHSVIS Executor")

	dconfig := exec.DriverConfig{
		Executor: newVISGHSExecutor(),
	}

	driver, err := exec.NewMesosExecutorDriver(dconfig)
	if err != nil {
		fmt.Println("Unable to create a ExecutorDriver ", err.Error())
	}

	_, err = driver.Start()
	if err != nil {
		fmt.Println("Got error:", err)
		return
	}
	fmt.Println("Executor process has started and running.")
	_, err = driver.Join()
	if err != nil {
		fmt.Println("driver failed:", err)
	}
	fmt.Println("executor terminating")

}
