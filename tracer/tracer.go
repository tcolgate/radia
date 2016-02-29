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

	pb "github.com/tcolgate/vonq/tracer/internal/proto"
)

type MessageDir pb.EdgeMessageRequest_Dir

var (
	DefaultTracer *Tracer

	EMDirIN    MessageDir = MessageDir(pb.EdgeMessageRequest_IN)
	EMDirOUT   MessageDir = MessageDir(pb.EdgeMessageRequest_OUT)
	EMDirQUEUE MessageDir = MessageDir(pb.EdgeMessageRequest_QUEUE)
)

func (m MessageDir) String() string {
	switch m {
	case EMDirIN:
		return "IN"
	case EMDirOUT:
		return "OUT"
	case EMDirQUEUE:
		return "QUEUE"
	default:
		return "UNKNOWN"
	}
}

//go:generate protoc -I $GOPATH/src:. --js_out=assets/ internal/proto/tracer.proto
type traceDisplay interface {
	Log(ts time.Time, id, s string)                                          // Log a plain text message
	NodeUpdate(ts time.Time, id, s string)                                   // Log change of state of a node
	EdgeUpdate(ts time.Time, id, edgeName, s string)                         // Log change of state of an edge
	EdgeMessage(ts time.Time, id, edgeName string, d MessageDir, str string) // Log send/recv of a message
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

func EdgeUpdate(id, edgeId, s string) {
	DefaultTracer.EdgeUpdate(id, edgeId, s)
}

func EdgeMessage(id, edgeId string, dir MessageDir, s string) {
	DefaultTracer.EdgeMessage(id, edgeId, dir, s)
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

func (t *Tracer) EdgeUpdate(id, edgeId, s string) {
	if t != nil {
		t.td.EdgeUpdate(time.Now(), id, edgeId, s)
	}
}

func (t *Tracer) EdgeMessage(id, edgeId string, dir MessageDir, str string) {
	if t != nil {
		t.td.EdgeMessage(time.Now(), id, edgeId, dir, str)
	}
}
