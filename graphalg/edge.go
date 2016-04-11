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
	"encoding/json"
	"fmt"
	"sort"

	"github.com/tcolgate/radia/graph"
	"github.com/tcolgate/radia/tracer"
)

type Edge struct {
	Weight   Weight
	Disabled bool

	local  *Node
	remote *Node

	SenderReciever
}

type Edges []*Edge

func (es Edges) String() string {
	str := "(Edges "
	for i, e := range es {
		str = str + fmt.Sprintf("(%v: %v)", i, e)
	}
	str += ")"
	return str
}

func (e Edge) String() string {
	return fmt.Sprintf("E(->%v:%v)", e.remote.ID, e.Weight)
}

func (e Edges) MinEdge() int {
	if len(e) == 0 {
		return -1
	}
	e.SortByMinEdge()
	return 0
}

func (e Edges) SortByMinEdge() {
	sort.Sort(ByMinEdge(e))
}

// ByMinEdge - implements sort by minimum edge
type ByMinEdge []*Edge

func (e ByMinEdge) Len() int {
	return len(e)
}

func (e ByMinEdge) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
func (e ByMinEdge) Less(i, j int) bool {
	return e[i].Weight.Less(e[j].Weight)
}

func NewEdge(f SenderRecieverMaker) (*Edge, *Edge) {
	c1, c2 := f()
	return &Edge{SenderReciever: c1}, &Edge{SenderReciever: c2}
}

func (e *Edge) Recieve() (interface{}, error) {
	i, err := e.SenderReciever.Recieve()
	if e.local.Tracer != nil {
		str, _ := json.Marshal(i)
		e.EdgeMessage(string(str), tracer.EMDirIN)
	}

	return i, err
}

func (e *Edge) Send(m MessageMarshaler) {
	e.SenderReciever.Send(m)
	if e.local.Tracer != nil {
		str, _ := json.Marshal(m)
		e.EdgeMessage(string(str), tracer.EMDirOUT)
	}
}

func (e *Edge) EdgeUpdate() {
	bs, _ := json.Marshal(e.Weight)
	sbs, _ := json.Marshal(e)
	e.local.Tracer.EdgeUpdate(graph.GraphID{}, graph.AlgorithmID{}, string(e.local.ID), string(bs), string(sbs))
}

func (e *Edge) EdgeMessage(s string, d tracer.MessageDir) {
	bs, _ := json.Marshal(e.Weight)
	e.local.Tracer.EdgeMessage(graph.GraphID{}, graph.AlgorithmID{}, string(e.local.ID), string(bs), d, s)
}
