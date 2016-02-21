package tracer

import (
	"io"
	"log"
)

func NewLogDisplay(w io.Writer) traceDisplay {
	return &logDisplay{log.New(w, "vonq: ", 0)}
}

type logDisplay struct {
	*log.Logger
}

func (l *logDisplay) Log(s string) {
	l.Println(s)
}

func (l *logDisplay) NodeUpdate() {
}

func (l *logDisplay) EdgeUpdate() {
}

func (l *logDisplay) EdgeMessage(str string) {
	l.Println(str)
}
