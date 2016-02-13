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

package graphalg

func Run(a Algorithm) {
	msgQueue := []Message{}
	ms := make(chan Message)
	defer close(ms)

	a.Edges().SortByMinEdge()
	defer func() {
		if a.OnDone() != nil {
			a.OnDone()
		}
	}()

	for ei, e := range a.Edges() {
		go func(e *Edge, ei int) {
			for {
				pb := e.Recieve()
				ms <- Message{ei, pb}
			}
		}(e, ei)
	}

	for nm := range ms {
		delayed := msgQueue
		msgQueue = []Message{}
		a.Dispatch(nm.Edge, nm.Data)
		for _, om := range delayed {
			a.Dispatch(om.Edge, om.Data)
		}
		if a.Done() {
			return
		}
	}
}
