package tracer

import "golang.org/x/net/context"

// tracerProxy
type tracerProxy struct {
	onward TraceServiceServer
}

func (s *tracerProxy) Log(ctx context.Context, message *LogRequest) (*LogResponse, error) {
	return &LogResponse{}, nil
}

func (s *tracerProxy) NodeUpdate(ctx context.Context, message *NodeUpdateRequest) (*NodeUpdateResponse, error) {
	return &NodeUpdateResponse{}, nil
}

func (s *tracerProxy) EdgeUpdate(ctx context.Context, message *EdgeUpdateRequest) (*EdgeUpdateResponse, error) {
	return &EdgeUpdateResponse{}, nil
}

func (s *tracerProxy) MessageUpdate(ctx context.Context, message *MessageUpdateRequest) (*MessageUpdateResponse, error) {
	return &MessageUpdateResponse{}, nil
}
