package tracer

import (
	"io"
	"log"
	"os"
)

var DefaultTracer Tracer

type Tracer struct {
	io.Writer
}

func init() {
	DefaultTracer = Tracer{os.Stdout}
}

func Print(v ...interface{}) {
	DefaultTracer.Print()
}

func Println(v ...interface{}) {
	DefaultTracer.Println()
}

func Printf(s string, v ...interface{}) {
	DefaultTracer.Printf(s, v...)
}

func NodeUpdate() {
}

func EdgeUpdate() {
}

func Message() {
}

func (*Tracer) Print(v ...interface{}) {
	log.Print(v...)
}

func (*Tracer) Println(v ...interface{}) {
	log.Println(v...)
}

func (*Tracer) Printf(s string, v ...interface{}) {
	log.Printf(s, v...)
}

func (*Tracer) NodeUpdate() {
}

func (*Tracer) EdgeUpdate() {
}

func (*Tracer) Message() {
}
