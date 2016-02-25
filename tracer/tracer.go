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

import "os"

var DefaultTracer *Tracer

//go:generate protoc -I $GOPATH/src:. --js_out=assets/ internal/proto/tracer.proto
type traceDisplay interface {
	Log(t int64, id, s string)
	NodeUpdate(t int64, id, s string)
	EdgeUpdate(t int64, id, edgeName, s string)
	EdgeMessage(t int64, id, edgeName, str string)
}

type Tracer struct {
	td traceDisplay
}

func init() {
	DefaultTracer = &Tracer{NewLogDisplay(os.Stdout)}
}

func Log(t int64, id, s string) {
	DefaultTracer.Log(t, id, s)
}

func NodeUpdate(t int64, id, s string) {
	DefaultTracer.NodeUpdate(t, id, s)
}

func EdgeUpdate(t int64, id, s string) {
	DefaultTracer.EdgeUpdate(t, id, s)
}

func EdgeMessage(t int64, id, edgeId, s string) {
	DefaultTracer.EdgeUpdate(t, id, s)
}

func (t *Tracer) Log(tm int64, id, s string) {
	if t != nil {
		t.td.Log(tm, id, s)
	}
}

func (t *Tracer) NodeUpdate(tm int64, id, s string) {
	if t != nil {
		t.td.NodeUpdate(tm, id, s)
	}
}

func (t *Tracer) EdgeUpdate(tm int64, id, s string) {
	if t != nil {
		t.td.NodeUpdate(tm, id, s)
	}
}

func (t *Tracer) EdgeMessage(tm int64, id, edgeId, str string) {
	if t != nil {
		t.td.EdgeMessage(tm, id, edgeId, str)
	}
}
