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

func (l *logDisplay) Log(t int64, id, s string) error {
	l.Printf("%v node(%v): %s", time.Unix(0, t), id, s)
	return nil
}

func (l *logDisplay) NodeUpdate(t int64, n, str string) error {
	l.Println(t, n, str)
	return nil
}

func (l *logDisplay) EdgeUpdate(t int64, n, en, str string) error {
	l.Println(t, n, en, str)
	return nil
}

func (l *logDisplay) EdgeMessage(t int64, n, en, str string) error {
	l.Println(t, n, en, str)
	return nil
}
