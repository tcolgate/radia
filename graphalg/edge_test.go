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

package graphalg

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func TestEdgeTestMessage1(t *testing.T) {
	n1 := Node{
		ID:     NodeID("n1"),
		Logger: log.New(os.Stdout, "node(n1) ", 0),
	}
	n2 := Node{
		ID:     NodeID("n2"),
		Logger: log.New(os.Stdout, "node(n2) ", 0),
	}
	Join(&n1, &n2, 1.0, MakeChanPair)

	if len(n1.Edges) != 1 {
		t.Fatalf("expected %v edges, got %v", 1, len(n1.Edges))
	}
	if len(n2.Edges) != 1 {
		t.Fatalf("expected %v edges, got %v", 1, len(n2.Edges))
	}

	sentm := ConnectMessage(0)
	go func() { n1.Edges[0].Send(sentm) }()
	gotm := n2.Edges[0].Recieve()

	gotm.Edge = nil
	if !reflect.DeepEqual(sentm, gotm) {
		t.Fatalf("expected %+v, got %+v", sentm, gotm)
	}
}

func TestEdgeTestMessage2(t *testing.T) {
	n1 := Node{
		ID:     NodeID("n1"),
		Logger: log.New(os.Stdout, "node(n1) ", 0),
	}
	n2 := Node{
		ID:     NodeID("n2"),
		Logger: log.New(os.Stdout, "node(n2) ", 0),
	}
	n3 := Node{
		ID:     NodeID("n3"),
		Logger: log.New(os.Stdout, "node(n3) ", 0),
	}
	Join(&n1, &n2, 1.0, MakeChanPair)
	Join(&n3, &n2, 1.0, MakeChanPair)

	if len(n1.Edges) != 1 {
		t.Fatalf("expected %v edges, got %v", 1, len(n1.Edges))
	}
	if len(n2.Edges) != 2 {
		t.Fatalf("expected %v edges, got %v", 1, len(n2.Edges))
	}

	if len(n3.Edges) != 1 {
		t.Fatalf("expected %v edges, got %v", 1, len(n2.Edges))
	}

	sentm := ConnectMessage(1)
	c := make(chan Message)
	go func() { c <- n1.Edges[0].Recieve() }()
	go func() { n3.Edges[0].Send(sentm) }()
	gotm := n2.Edges[1].Recieve()
	select {
	case <-c:
		t.Fatalf("expected recieve to block")
	default:
	}

	gotm.Edge = nil
	if !reflect.DeepEqual(sentm, gotm) {
		t.Fatalf("expected %+v, got %+v", sentm, gotm)
	}
}

func TestEdgesNextEmpty(t *testing.T) {

}
