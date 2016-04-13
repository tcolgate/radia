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
)

// grpcPair is a sender reciever using channels
type grpcPair struct {
	send chan<- Message
	recv <-chan Message
}

type grpcProxy struct {
	sync.Mutex
	smap map[NodeID]*grpc.ClientConn
	rmap map[NodeID]chan<- Message
}

// Subscribe returns a Reciever for the remote NodeID
// passed in.
func (p *grpcProxy) Subscribe(r NodeID) Reciever {
	return nil
}

func NewGRPCServer() MessageServiceServer {
	return &grpcProxy{}
}

func (*grpcProxy) SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error) {
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

// MakeGRPCPair is an edge sender/reciever built using a
// channel
func MakeGRPCPair() (SenderReciever, SenderReciever) {
	c1, c2 := make(chan Message), make(chan Message)
	return grpcPair{c1, c2}, grpcPair{c2, c1}
}
