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
	pb "github.com/tcolgate/vonq/tracer/internal/proto"
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

func (g *grpcClientDisplay) Log(t int64, id, s string) {
	r := pb.LogRequest{Time: t, NodeID: id, Message: s}
	g.TraceServiceClient.Log(context.Background(), &r)
}

func (g *grpcClientDisplay) NodeUpdate(t int64, id, s string) {
	r := pb.NodeUpdateRequest{}
	g.TraceServiceClient.NodeUpdate(context.Background(), &r)
}

func (g *grpcClientDisplay) EdgeUpdate(t int64, id, eid, s string) {
	r := pb.EdgeUpdateRequest{}
	g.TraceServiceClient.EdgeUpdate(context.Background(), &r)
}

func (g *grpcClientDisplay) EdgeMessage(t int64, id, eid, str string) {
	r := pb.EdgeMessageRequest{}
	g.TraceServiceClient.EdgeMessage(context.Background(), &r)
}

type grpcServerDisplay struct {
	o traceDisplay
}

func NewGRPCServer(onward traceDisplay) pb.TraceServiceServer {
	return &grpcServerDisplay{onward}
}

func (s *grpcServerDisplay) Log(ctx context.Context, r *pb.LogRequest) (*pb.LogResponse, error) {
	s.o.Log(r.Time, r.NodeID, r.Message)
	return &pb.LogResponse{}, nil
}

func (s *grpcServerDisplay) NodeUpdate(ctx context.Context, r *pb.NodeUpdateRequest) (*pb.NodeUpdateResponse, error) {
	s.o.NodeUpdate(r.Time, r.NodeID, r.Status)
	return &pb.NodeUpdateResponse{}, nil
}

func (s *grpcServerDisplay) EdgeUpdate(ctx context.Context, r *pb.EdgeUpdateRequest) (*pb.EdgeUpdateResponse, error) {
	s.o.EdgeUpdate(r.Time, r.NodeID, r.EdgeName, r.Status)
	return &pb.EdgeUpdateResponse{}, nil
}

func (s *grpcServerDisplay) EdgeMessage(ctx context.Context, r *pb.EdgeMessageRequest) (*pb.EdgeMessageResponse, error) {
	s.o.EdgeMessage(r.Time, r.NodeID, r.EdgeName, r.Message)
	return &pb.EdgeMessageResponse{}, nil
}
