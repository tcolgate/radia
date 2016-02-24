package tracer

import "os"

var DefaultTracer *Tracer

//go:generate protoc -I $GOPATH/src:. --js_out=assets/ internal/proto/tracer.proto
type traceDisplay interface {
	Log(t int64, id, s string) error
	NodeUpdate(t int64, id, s string) error
	EdgeUpdate(t int64, id, edgeName, s string) error
	EdgeMessage(t int64, id, edgeName, str string) error
}

type Tracer struct {
	td traceDisplay
}

func init() {
	DefaultTracer = &Tracer{NewLogDisplay(os.Stdout)}
}

func Log(t int64, id, s string) error {
	DefaultTracer.Log(t, id, s)
	return nil
}

func NodeUpdate(t int64, id, s string) error {
	DefaultTracer.NodeUpdate(t, id, s)
	return nil
}

func EdgeUpdate(t int64, id, s string) error {
	DefaultTracer.EdgeUpdate(t, id, s)
	return nil
}

func EdgeMessage(t int64, id, edgeId, s string) error {
	DefaultTracer.EdgeUpdate(t, id, s)
	return nil
}

func (t *Tracer) Log(tm int64, id, s string) error {
	if t != nil {
		t.td.Log(tm, id, s)
	}
	return nil
}

func (t *Tracer) NodeUpdate(tm int64, id, s string) error {
	if t != nil {
		t.td.NodeUpdate(tm, id, s)
	}
	return nil
}

func (t *Tracer) EdgeUpdate(tm int64, id, s string) error {
	if t != nil {
		t.td.NodeUpdate(tm, id, s)
	}
	return nil
}

func (t *Tracer) EdgeMessage(tm int64, id, edgeId, str string) error {
	if t != nil {
		t.td.EdgeMessage(tm, id, edgeId, str)
	}
	return nil
}
