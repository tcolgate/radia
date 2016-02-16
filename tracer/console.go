package tracer

import "golang.org/x/net/context"

type console struct{}

func (s *console) Log(ctx context.Context, message *LogRequest) (*LogResponse, error) {
	return &LogResponse{}, nil
}

func (s *console) NodeUpdate(ctx context.Context, message *NodeUpdateRequest) (*NodeUpdateResponse, error) {
	return &NodeUpdateResponse{}, nil
}

func (s *console) EdgeUpdate(ctx context.Context, message *EdgeUpdateRequest) (*EdgeUpdateResponse, error) {
	return &EdgeUpdateResponse{}, nil
}

func (s *console) MessageUpdate(ctx context.Context, message *MessageUpdateRequest) (*MessageUpdateResponse, error) {
	return &MessageUpdateResponse{}, nil
}
