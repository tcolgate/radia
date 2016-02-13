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
	"sort"
)

type NodeID string

type Node struct {
	Base
	ID NodeID

	edges Edges

	msgQueue []Message

	OnDone func()

	*log.Logger
}

func (n *Node) Edges() Edges {
	return n.edges
}

func (n *Node) MinEdge() int {
	return n.edges.MinEdge()
}

func (n *Node) Edge(j int) *Edge {
	return n.edges[j]
}

func Join(n1 *Node, n2 *Node, w float64, f SenderRecieverMaker) {
	n1.Join(n2, w, f)
}

func (n1 *Node) Join(n2 *Node, w float64, f SenderRecieverMaker) {
	ids := []string{
		string(n1.ID),
		string(n2.ID),
	}
	sort.Strings(ids)

	e1, e2 := NewEdge(f)

	e1.Weight.Cost = w
	e2.Weight.Cost = w

	e1.Weight.LsnID = ids[0]
	e2.Weight.LsnID = ids[0]
	e1.Weight.MsnID = ids[1]
	e2.Weight.MsnID = ids[1]

	e1.local, e1.remote = n1, n2 // mostly for debugging
	e2.local, e2.remote = n2, n1

	n1.edges = append(n1.edges, e1)
	n2.edges = append(n2.edges, e2)
}

// Send - send a message to the specified
func (n *Node) Send(e int, d []byte) {
	n.edges[e].Send(d)
}

// Queue - re-queue a message to the internal queue
func (n *Node) Queue(e int, d []byte) {
	n.msgQueue = append(n.msgQueue, Message{e, d})
}
