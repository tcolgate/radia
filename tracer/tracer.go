package tracer

import "os"

var DefaultTracer *Tracer

type traceDisplay interface {
	Log(s string)
	NodeUpdate()
	EdgeUpdate()
	EdgeMessage(str string)
}

type Tracer struct {
	td traceDisplay
}

func init() {
	DefaultTracer = &Tracer{NewLogDisplay(os.Stdout)}
}

func Log(s string) {
	DefaultTracer.Log(s)
}

func NodeUpdate() {
	DefaultTracer.NodeUpdate()
}

func EdgeUpdate() {
	DefaultTracer.EdgeUpdate()
}

func (t *Tracer) Log(s string) {
	if t != nil {
		t.td.Log(s)
	}
}

func (t *Tracer) NodeUpdate() {
	if t != nil {
		t.td.NodeUpdate()
	}
}

func (t *Tracer) EdgeUpdate() {
	if t != nil {
		t.td.NodeUpdate()
	}
}

func (t *Tracer) EdgeMessage(str string) {
	if t != nil {
		t.td.EdgeMessage(str)
	}
}
