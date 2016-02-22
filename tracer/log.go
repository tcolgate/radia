package tracer

import (
	"io"
	"log"
	"time"
)

func NewLogDisplay(w io.Writer) traceDisplay {
	return &logDisplay{log.New(w, "vonq: ", 0)}
}

type logDisplay struct {
	*log.Logger
}

func (l *logDisplay) Log(t int64, id, s string) {
	l.Printf("%v node(%v): %s", time.Unix(0, t), id, s)
}

func (l *logDisplay) NodeUpdate(t int64, n, str string) {
	l.Println(t, n, str)
}

func (l *logDisplay) EdgeUpdate(t int64, n, en, str string) {
	l.Println(t, n, en, str)
}

func (l *logDisplay) EdgeMessage(t int64, n, en, str string) {
	l.Println(t, n, en, str)
}
