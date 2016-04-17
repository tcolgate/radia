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
	"net"
	"testing"

	"google.golang.org/grpc"
)

var cat = map[NodeID]string{}

func TestRPCPairTestMessage1(t *testing.T) {
	RegisterMessage(TestMessage{})

	l1, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err.Error())
	}

	n1 := Node{
		ID: NodeID(l1.Addr().String()),
	}
	cat[n1.ID] = l1.Addr().String()

	s1 := grpc.NewServer()
	sp1 := NewGRPCServer()
	RegisterMessageServiceServer(s1, sp1)

	l2, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err.Error())
	}
	n2 := Node{
		ID: NodeID(l2.Addr().String()),
	}
	cat[n2.ID] = l2.Addr().String()
	s2 := grpc.NewServer()
	sp2 := NewGRPCServer()
	RegisterMessageServiceServer(s2, sp2)

	/*
		Join(&n1, &n2, 1.0, MakeGRPCPair)

		if len(n1.Edges()) != 1 {
			t.Fatalf("expected %v edges, got %v", 1, len(n1.Edges()))
		}
		if len(n2.Edges()) != 1 {
			t.Fatalf("expected %v edges, got %v", 1, len(n2.Edges()))
		}

		sentm := &TestMessage{1}
		go func() { n1.Edges()[0].Send(sentm) }()

		goti, err := n2.Edges()[0].Recieve()
		if err != nil {
			t.Fatalf("error recieving, %v", err)
		}

		gotm, ok := goti.(*TestMessage)

		if !ok {
			log.Fatalf("failure type switching message")
		}

		//	gotm.Edge = nil
		if !reflect.DeepEqual(sentm, gotm) {
			t.Fatalf("expected %+v, got %+v", sentm, gotm)
		}
	*/
}
