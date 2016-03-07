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

	"github.com/tcolgate/vonq/graph"
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

type traceDisplay interface {
	Log(ts time.Time, gid graph.GraphID, aid graph.AlgorithmID, id, s string)                                          // Log a plain text message
	NodeUpdate(ts time.Time, gid graph.GraphID, aid graph.AlgorithmID, id, s string)                                   // Log change of state of a node
	EdgeUpdate(ts time.Time, gid graph.GraphID, aid graph.AlgorithmID, id, edgeName, s string)                         // Log change of state of an edge
	EdgeMessage(ts time.Time, gid graph.GraphID, aid graph.AlgorithmID, id, edgeName string, d MessageDir, str string) // Log send/recv of a message
}

type Tracer struct {
	td traceDisplay
}

func init() {
	DefaultTracer = &Tracer{NewLogDisplay(os.Stdout)}
}

func Log(gid graph.GraphID, aid graph.AlgorithmID, id, s string) {
	DefaultTracer.Log(gid, aid, id, s)
}

func NodeUpdate(gid graph.GraphID, aid graph.AlgorithmID, id, s string) {
	DefaultTracer.NodeUpdate(gid, aid, id, s)
}

func EdgeUpdate(gid graph.GraphID, aid graph.AlgorithmID, id, edgeId, s string) {
	DefaultTracer.EdgeUpdate(gid, aid, id, edgeId, s)
}

func EdgeMessage(gid graph.GraphID, aid graph.AlgorithmID, id, edgeId string, dir MessageDir, s string) {
	DefaultTracer.EdgeMessage(gid, aid, id, edgeId, dir, s)
}

func (t *Tracer) Log(gid graph.GraphID, aid graph.AlgorithmID, id, s string) {
	if t != nil {
		t.td.Log(time.Now(), gid, aid, id, s)
	}
}

func (t *Tracer) NodeUpdate(gid graph.GraphID, aid graph.AlgorithmID, id, s string) {
	if t != nil {
		t.td.NodeUpdate(time.Now(), gid, aid, id, s)
	}
}

func (t *Tracer) EdgeUpdate(gid graph.GraphID, aid graph.AlgorithmID, id, edgeId, s string) {
	if t != nil {
		t.td.EdgeUpdate(time.Now(), gid, aid, id, edgeId, s)
	}
}

func (t *Tracer) EdgeMessage(gid graph.GraphID, aid graph.AlgorithmID, id, edgeId string, dir MessageDir, str string) {
	if t != nil {
		t.td.EdgeMessage(time.Now(), gid, aid, id, edgeId, dir, str)
	}
}

func New(td traceDisplay) *Tracer {
	return &Tracer{td}
}
