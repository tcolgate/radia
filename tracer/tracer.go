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

package tracer

import (
	"os"
	"time"
)

var DefaultTracer *Tracer

//go:generate protoc -I $GOPATH/src:. --js_out=assets/ internal/proto/tracer.proto
type traceDisplay interface {
	Log(ts time.Time, id, s string)                     // Log a plain text message
	NodeUpdate(ts time.Time, id, s string)              // Log change of state of a node
	EdgeUpdate(ts time.Time, id, edgeName, s string)    // Log change of state of an edge
	EdgeMessage(ts time.Time, id, edgeName, str string) // Log send/recv of a message
}

type Tracer struct {
	td traceDisplay
}

func init() {
	DefaultTracer = &Tracer{NewLogDisplay(os.Stdout)}
}

func Log(id, s string) {
	DefaultTracer.Log(id, s)
}

func NodeUpdate(id, s string) {
	DefaultTracer.NodeUpdate(id, s)
}

func EdgeUpdate(id, s string) {
	DefaultTracer.EdgeUpdate(id, s)
}

func EdgeMessage(id, edgeId, s string) {
	DefaultTracer.EdgeUpdate(id, s)
}

func (t *Tracer) Log(id, s string) {
	if t != nil {
		t.td.Log(time.Now(), id, s)
	}
}

func (t *Tracer) NodeUpdate(id, s string) {
	if t != nil {
		t.td.NodeUpdate(time.Now(), id, s)
	}
}

func (t *Tracer) EdgeUpdate(id, s string) {
	if t != nil {
		t.td.NodeUpdate(time.Now(), id, s)
	}
}

func (t *Tracer) EdgeMessage(id, edgeId, str string) {
	if t != nil {
		t.td.EdgeMessage(time.Now(), id, edgeId, str)
	}
}
