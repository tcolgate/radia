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
