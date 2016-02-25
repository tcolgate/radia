// Copyright (c) 2016 Tristan Colgate-McFarlane
//
// This file is part of vonq.
//
// vonq is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// vonq is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with vonq.  If not, see <http://www.gnu.org/licenses/>.

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
