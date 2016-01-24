package udpping

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/tcolgate/vonq/internal/reporter"
	pb "github.com/tcolgate/vonq/probe/udpping/proto"
)

type client struct {
	addr net.UDPAddr
	key  []byte
	r    reporter.Reporter
}

func (cl *client) sendPing(c *net.UDPConn, count uint32) error {
	thing := pb.PingRequest{}
	tns := uint64(time.Now().UnixNano())
	thing.Time = &tns
	thing.Count = &count
	bs, err := proto.Marshal(&thing)
	if err != nil {
		panic(err)
	}

	mac := genMAC(bs, cl.key)

	b := append(mac, bs...)
	_, e := c.Write(b)

	log.Println(count)

	return e
}

func (c *client) run(daddr net.UDPAddr) {
	var count uint32
	s, err := net.DialUDP("udp", &c.addr, &daddr)
	if err != nil {
		os.Exit(1)
	}
	go c.getReplies(s)

	for {
		c.sendPing(s, count)
		time.Sleep(time.Second)
		count++
	}

	s.Close()
}

func (c *client) getReplies(so *net.UDPConn) {
	for {
		b := make([]byte, 128)
		i, sa, e := so.ReadFromUDP(b)
		log.Println(i, e)

		tns := uint64(time.Now().UnixNano())
		macLen := 256 / 8
		mac := b[:macLen]
		message := b[macLen:i]

		if len(mac) != macLen {
			log.Println("Client: Short HMAC")
			continue
		}
		if !checkMAC(message, mac, c.key) {
			log.Println("Client: Bad HMAC ", macLen)
			continue
		}

		log.Println("Client: Good HMAC ", macLen)

		rep := pb.PingReply{}
		proto.Unmarshal(message, &rep)
		go c.processReply(&rep, tns, sa)
	}
}

func (c *client) processReply(r *pb.PingReply, timeIn uint64, sa *net.UDPAddr) {
	t1 := time.Unix(0, int64(*r.TimeSent))
	t2 := time.Unix(0, int64(*r.TimeIn))
	t3 := time.Unix(0, int64(*r.TimeOut))
	t4 := time.Unix(0, int64(timeIn))

	dout := t2.Sub(t1)
	dsrv := t3.Sub(t2)
	dback := t4.Sub(t3)
	rtt := t4.Sub(t1) - dsrv

	log.Println("ADDR: ", sa.IP)
	log.Println("DOUT: ", dout)
	log.Println("DSRV: ", dsrv)
	log.Println("DBACK: ", dback)
	log.Println("RTT: ", rtt)
	return
}
