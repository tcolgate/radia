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
	"encoding/json"
	"fmt"
	"sort"
)

import (
	"github.com/tcolgate/vonq/graph"
	"github.com/tcolgate/vonq/tracer"
)

type NodeID string

type QueuedMessage struct {
	e int
	m interface{}
}

type Node struct {
	Base
	ID NodeID

	edges Edges

	msgQueue []QueuedMessage

	*tracer.Tracer
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

	e1.Weight.Lsn = ids[0]
	e2.Weight.Lsn = ids[0]
	e1.Weight.Msn = ids[1]
	e2.Weight.Msn = ids[1]

	e1.local, e1.remote = n1, n2 // mostly for debugging
	e2.local, e2.remote = n2, n1

	n1.edges = append(n1.edges, e1)
	n2.edges = append(n2.edges, e2)
}

// Send - send a message to the specified
func (n *Node) Send(e int, d MessageMarshaler) {
	n.edges[e].Send(d)
}

// Queue - re-queue a message to the internal queue
func (n *Node) Queue(e int, d interface{}) {
	if n.Tracer != nil {
		str, _ := json.Marshal(d)
		n.Edges()[e].EdgeMessage(string(str), tracer.EMDirQUEUE)
	}
	n.msgQueue = append(n.msgQueue, QueuedMessage{e, d})
}

func (n *Node) Queued() []QueuedMessage {
	return n.msgQueue
}

func (n *Node) ClearQueue() {
	n.msgQueue = []QueuedMessage{}
}

func (n *Node) Log(s string) {
	n.Tracer.Log(graph.GraphID{}, graph.AlgorithmID{}, string(n.ID), s)
}

func (n *Node) Print(v ...interface{}) {
	n.Log(fmt.Sprint(v...))
}

func (n *Node) Println(v ...interface{}) {
	n.Log(fmt.Sprintln(v...))
}

func (n *Node) Printf(str string, v ...interface{}) {
	n.Log(fmt.Sprintf(str, v...))
}

func (n *Node) NodeUpdate(s string) {
	n.Tracer.NodeUpdate(graph.GraphID{}, graph.AlgorithmID{}, string(n.ID), s)
}
