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

func (g *grpcClientDisplay) Log(s string) {
	r := pb.LogRequest{Message: s}
	g.TraceServiceClient.Log(context.Background(), &r)
}

func (g *grpcClientDisplay) NodeUpdate() {
	r := pb.NodeUpdateRequest{}
	g.TraceServiceClient.NodeUpdate(context.Background(), &r)
}

func (g *grpcClientDisplay) EdgeUpdate() {
	r := pb.EdgeUpdateRequest{}
	g.TraceServiceClient.EdgeUpdate(context.Background(), &r)
}

func (g *grpcClientDisplay) EdgeMessage(str string) {
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
	s.o.Log(r.Message)
	return &pb.LogResponse{}, nil
}

func (s *grpcServerDisplay) NodeUpdate(ctx context.Context, r *pb.NodeUpdateRequest) (*pb.NodeUpdateResponse, error) {
	s.o.NodeUpdate()
	return &pb.NodeUpdateResponse{}, nil
}

func (s *grpcServerDisplay) EdgeUpdate(ctx context.Context, r *pb.EdgeUpdateRequest) (*pb.EdgeUpdateResponse, error) {
	s.o.EdgeUpdate()
	return &pb.EdgeUpdateResponse{}, nil
}

func (s *grpcServerDisplay) EdgeMessage(ctx context.Context, r *pb.EdgeMessageRequest) (*pb.EdgeMessageResponse, error) {
	s.o.EdgeMessage(r.String())
	return &pb.EdgeMessageResponse{}, nil
}
