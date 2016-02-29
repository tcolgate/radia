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

func Run(a NodeAlgorithm, onDone func()) {
	ms := make(chan QueuedMessage)
	defer close(ms)

	a.Edges().SortByMinEdge()
	defer func() {
		if onDone != nil {
			onDone()
		}
	}()

	for ei, e := range a.Edges() {
		go func(e *Edge, ei int) {
			for {
				pb := e.Recieve()
				ms <- QueuedMessage{ei, pb}
			}
		}(e, ei)
	}

	for nm := range ms {
		delayed := a.Queued()
		a.ClearQueue()

		a.Dispatch(nm.e, nm.m)
		for _, om := range delayed {
			a.Dispatch(om.e, om.m)
		}
		if a.Done() {
			return
		}
	}
}
