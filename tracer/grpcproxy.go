package tracer

import pb "github.com/tcolgate/vonq/tracer/internal/proto"
import "golang.org/x/net/context"

// tracerProxy
type tracerProxy struct {
	onward pb.TraceServiceServer
}

func (s *tracerProxy) Log(ctx context.Context, message *pb.LogRequest) (*pb.LogResponse, error) {
	return &pb.LogResponse{}, nil
}

func (s *tracerProxy) NodeUpdate(ctx context.Context, message *pb.NodeUpdateRequest) (*pb.NodeUpdateResponse, error) {
	return &pb.NodeUpdateResponse{}, nil
}

func (s *tracerProxy) EdgeUpdate(ctx context.Context, message *pb.EdgeUpdateRequest) (*pb.EdgeUpdateResponse, error) {
	return &pb.EdgeUpdateResponse{}, nil
}

func (s *tracerProxy) MessageUpdate(ctx context.Context, message *pb.MessageUpdateRequest) (*pb.MessageUpdateResponse, error) {
	return &pb.MessageUpdateResponse{}, nil
}
