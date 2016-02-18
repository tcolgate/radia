package tracer

import (
	"io"
	"os"
)

var DefaultTracer Tracer

type Tracer struct {
	io.Writer
}

func init() {
	DefaultTracer = Tracer{os.Stdout}
}

func Log() {
	DefaultTracer.Log()
}

func Print() {
	DefaultTracer.Print()
}

func Println() {
	DefaultTracer.Println()
}

func Printf() {
	DefaultTracer.Printf()
}

func NodeUpdate() {
}

func EdgeUpdate() {
}

func Message() {
}

func (Tracer) Log() {
}

func (Tracer) Print() {
}

func (Tracer) Println() {
}

func (Tracer) Printf() {
}

func (Tracer) NodeUpdate() {
}

func (Tracer) EdgeUpdate() {
}

func (Tracer) Message() {
}
