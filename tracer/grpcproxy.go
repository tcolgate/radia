package tracer

import (
	"io"
	"log"

	pb "github.com/tcolgate/vonq/tracer/internal/proto"
	"google.golang.org/grpc"
)
import "golang.org/x/net/context"

type grpcClientDisplay struct {
	pb.TraceServiceClient
	lc  pb.TraceService_LogClient
	nuc pb.TraceService_NodeUpdateClient
	euc pb.TraceService_EdgeUpdateClient
	emc pb.TraceService_EdgeMessageClient
}

func NewGRPCDisplayClient(addr string, os ...grpc.DialOption) (traceDisplay, error) {
	conn, err := grpc.Dial(addr, os...)
	p := pb.NewTraceServiceClient(conn)
	return &grpcClientDisplay{
		TraceServiceClient: p,
	}, err
}

func (g *grpcClientDisplay) Log(t int64, id, s string) (err error) {
	if g.lc == nil {
		g.lc, err = g.TraceServiceClient.Log(context.Background())
		if err != nil {
			return err
		}
	}
	return g.lc.Send(&pb.LogRequest{Time: t, NodeID: id, Message: s})
}

func (g *grpcClientDisplay) NodeUpdate(t int64, id, s string) (err error) {
	if g.nuc == nil {
		g.nuc, err = g.TraceServiceClient.NodeUpdate(context.Background())
		if err != nil {
			return err
		}
	}
	return g.nuc.Send(&pb.NodeUpdateRequest{})
}

func (g *grpcClientDisplay) EdgeUpdate(t int64, id, eid, s string) (err error) {
	if g.euc == nil {
		g.euc, err = g.TraceServiceClient.EdgeUpdate(context.Background())
		if err != nil {
			return err
		}
	}
	return g.euc.Send(&pb.EdgeUpdateRequest{})
}

func (g *grpcClientDisplay) EdgeMessage(t int64, id, eid, str string) (err error) {
	if g.emc == nil {
		g.emc, err = g.TraceServiceClient.EdgeMessage(context.Background())
		if err != nil {
			return err
		}
	}
	return g.emc.Send(&pb.EdgeMessageRequest{})
}

type grpcServerDisplay struct {
	o traceDisplay
}

func NewGRPCServer(onward traceDisplay) pb.TraceServiceServer {
	return &grpcServerDisplay{o: onward}
}

func (s *grpcServerDisplay) Log(rs pb.TraceService_LogServer) error {
	for {
		log.Println("HERE")
		r, err := rs.Recv()
		if err == io.EOF {
			return rs.SendAndClose(&pb.LogResponse{})
		}
		if err != nil {
			return err
		}
		s.o.Log(r.Time, r.NodeID, r.Message)
	}
}

func (s *grpcServerDisplay) NodeUpdate(rs pb.TraceService_NodeUpdateServer) error {
	for {
		r, err := rs.Recv()
		log.Println(r)
		if err == io.EOF {
			return rs.SendAndClose(&pb.NodeUpdateResponse{})
		}
		if err != nil {
			return err
		}
		s.o.NodeUpdate(r.Time, r.NodeID, r.Status)
	}
}

func (s *grpcServerDisplay) EdgeUpdate(rs pb.TraceService_EdgeUpdateServer) error {
	for {
		r, err := rs.Recv()
		log.Println(r)
		if err == io.EOF {
			return rs.SendAndClose(&pb.EdgeUpdateResponse{})
		}
		if err != nil {
			return err
		}
		s.o.EdgeUpdate(r.Time, r.NodeID, r.EdgeName, r.Status)
	}
	return nil
}

func (s *grpcServerDisplay) EdgeMessage(rs pb.TraceService_EdgeMessageServer) error {
	for {
		r, err := rs.Recv()
		log.Println(r)
		if err == io.EOF {
			return rs.SendAndClose(&pb.EdgeMessageResponse{})
		}
		if err != nil {
			return err
		}
		s.o.EdgeMessage(r.Time, r.NodeID, r.EdgeName, r.Message)
	}
	return nil
}
