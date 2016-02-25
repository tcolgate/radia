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
