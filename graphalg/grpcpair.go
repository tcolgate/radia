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

package graphalg

import (
	"sync"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	google_protobuf "github.com/golang/protobuf/ptypes/any"
	"github.com/tcolgate/radia/graph"
)

// grpcPair is a sender reciever using channels
type grpcPair struct {
	g    graph.GraphID
	a    graph.AlgorithmID
	send chan<- Message
	recv <-chan Message
}

type endAddress struct {
	g graph.GraphID
	a graph.AlgorithmID
}

type grpcProxy struct {
	sync.Mutex
	smap map[NodeID]*grpc.ClientConn
	rmap map[endAddress]chan<- Message
}

// subscribe returns a Reciever for the remote NodeID
// passed in.
func (p *grpcProxy) subscribe(g graph.GraphID, a graph.AlgorithmID) chan<- Message {
	p.Lock()
	defer p.Unlock()
	if p.rmap == nil {
		p.rmap = make(map[endAddress]chan<- Message)
	}

	if _, ok := p.rmap[endAddress{g, a}]; ok {
		panic("already subscribed")
	}

	c := make(chan Message)

	p.rmap[endAddress{g, a}] = c
	return c
}

func NewGRPCServer() MessageServiceServer {
	return &grpcProxy{}
}

// SendMessage implements the SendMessage RPC - this is going to be very difficult
// we need to look at errors, cancellation, possibly duplicate message detecton
// integrity. We'll start off hugely naive.
func (p *grpcProxy) SendMessage(ctx context.Context, r *SendMessageRequest) (*SendMessageResponse, error) {
	if c, ok := p.rmap[endAddress{*r.Gid, *r.Aid}]; ok {
		return nil, nil
	} else {
		c <- *r.Msg
		return nil, nil
	}
}

func (*grpcProxy) EdgeWeight(context.Context, *EdgeWeightRequest) (*EdgeWeightResponse, error) {
	return nil, nil
}

func (p grpcPair) Send(m MessageMarshaler) {
	bs, url := m.MarshalMessage()

	p.send <- Message{
		Payload: &google_protobuf.Any{
			TypeUrl: url,
			Value:   bs,
		},
	}
}

func (p grpcPair) Recieve() (interface{}, error) {
	m := <-p.recv

	return unmarshalAny(m.Payload)
}

func (p grpcPair) Close() {
	close(p.send)
}

// MakeGRPCEdge is an edge sender/reciever built using a
// channel
func MakeGRPCEdge(g graph.GraphID, a graph.AlgorithmID) SenderReciever {
	return grpcPair{g: g, a: a}
}
