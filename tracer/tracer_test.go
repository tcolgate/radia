package tracer

import (
	"bytes"
	"log"
	"net"
	"testing"

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
	c.Log("Something")
	exp := "vonq: Something\n"
	if got := out.String(); got != exp {
		t.Fatalf("expected %v got: %v", exp, got)
	}
}
