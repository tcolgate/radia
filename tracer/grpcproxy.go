// Copyright (c) 2016 Tristan Colgate-McFarlane
//
// This file is part of radia.
//
// radia is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// radia is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with radia.  If not, see <http://www.gnu.org/licenses/>.

package tracer

import (
	"time"

	"github.com/tcolgate/radia/graph"
	pb "github.com/tcolgate/radia/tracer/internal/proto"
	"google.golang.org/grpc"
)
import "golang.org/x/net/context"

type grpcClientDisplay struct {
	pb.TraceServiceClient
}

func NewGRPCDisplayClient(addr string, os ...grpc.DialOption) (traceDisplay, error) {
	conn, err := grpc.Dial(addr, os...)
	return &grpcClientDisplay{
		pb.NewTraceServiceClient(conn),
	}, err
}

func (g *grpcClientDisplay) Log(t time.Time, gid graph.GraphID, aid graph.AlgorithmID, id, s string) {
	r := pb.LogRequest{Time: t.UnixNano(), NodeID: id, Message: s, Gid: &gid, Aid: &aid}
	g.TraceServiceClient.Log(context.Background(), &r)
}

func (g *grpcClientDisplay) NodeUpdate(t time.Time, gid graph.GraphID, aid graph.AlgorithmID, id, s string) {
	r := pb.NodeUpdateRequest{}
	g.TraceServiceClient.NodeUpdate(context.Background(), &r)
}

func (g *grpcClientDisplay) EdgeUpdate(t time.Time, gid graph.GraphID, aid graph.AlgorithmID, id, eid, s string) {
	r := pb.EdgeUpdateRequest{}
	g.TraceServiceClient.EdgeUpdate(context.Background(), &r)
}

func (g *grpcClientDisplay) EdgeMessage(t time.Time, gid graph.GraphID, aid graph.AlgorithmID, id, eid string, dir MessageDir, str string) {
	r := pb.EdgeMessageRequest{}
	g.TraceServiceClient.EdgeMessage(context.Background(), &r)
}

type grpcServerDisplay struct {
	o traceDisplay
}

func NewGRPCServer(onward traceDisplay) pb.TraceServiceServer {
	return &grpcServerDisplay{onward}
}

func RegisterTraceServiceServer(gs *grpc.Server, server pb.TraceServiceServer) {
	pb.RegisterTraceServiceServer(gs, server)
}

func (s *grpcServerDisplay) Log(ctx context.Context, r *pb.LogRequest) (*pb.LogResponse, error) {
	s.o.Log(time.Unix(9, r.Time), *r.Gid, *r.Aid, r.NodeID, r.Message)
	return &pb.LogResponse{}, nil
}

func (s *grpcServerDisplay) NodeUpdate(ctx context.Context, r *pb.NodeUpdateRequest) (*pb.NodeUpdateResponse, error) {
	s.o.NodeUpdate(time.Unix(0, r.Time), *r.Gid, *r.Aid, r.NodeID, r.Status)
	return &pb.NodeUpdateResponse{}, nil
}

func (s *grpcServerDisplay) EdgeUpdate(ctx context.Context, r *pb.EdgeUpdateRequest) (*pb.EdgeUpdateResponse, error) {
	s.o.EdgeUpdate(time.Unix(0, r.Time), *r.Gid, *r.Aid, r.NodeID, r.EdgeName, r.Status)
	return &pb.EdgeUpdateResponse{}, nil
}

func (s *grpcServerDisplay) EdgeMessage(ctx context.Context, r *pb.EdgeMessageRequest) (*pb.EdgeMessageResponse, error) {
	s.o.EdgeMessage(time.Unix(0, r.Time), *r.Gid, *r.Aid, r.NodeID, r.EdgeName, MessageDir(r.Direction), r.Message)
	return &pb.EdgeMessageResponse{}, nil
}
