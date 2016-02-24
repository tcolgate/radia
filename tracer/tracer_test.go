package tracer

import (
	"bytes"
	"log"
	"net"
	"testing"
	"time"

	pb "github.com/tcolgate/vonq/tracer/internal/proto"
	"google.golang.org/grpc"
)

func TestTracerServer1(t *testing.T) {
	out := bytes.Buffer{}
	outer := NewLogDisplay(&out)
	server := NewGRPCServer(outer)

	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTraceServiceServer(grpcServer, server)
	go grpcServer.Serve(lis)

	log.Println(lis.Addr())
	c, err := NewGRPCDisplayClient(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	err = c.Log(1, "2", "Something")
	log.Println("HERE AGAIN")
	time.Sleep(1 * time.Second) // This isn't going to work
	if err != nil {
		log.Fatalf("failed to call client: %v", err)
	}
	exp := "vonq: Something\n"
	if got := out.String(); got != exp {
		t.Fatalf("expected %v got: %v", exp, got)
	}
}
