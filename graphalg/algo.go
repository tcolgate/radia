// Copyright (c) 2016 Tristan Colgate-McFarlane
//
// This file is part of radia.
//
// radia is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// radia is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with radia.  If not, see <http://www.gnu.org/licenses/>.

package graphalg

type NodeAlgorithm interface {
	Dispatcher
	Doner
}

type Queuer interface {
	Queue(edge int, m interface{})
}

type Dispatcher interface {
	Dispatch(edge int, m interface{})
	Edges() Edges
	Queued() []QueuedMessage
	ClearQueue()
}

type Doner interface {
	Done() bool
}

type Base struct {
	IsDone bool
	OnDone func()
}

func (b *Base) Done() bool {
	return b.IsDone
}

func (b *Base) SetDone(done bool) {
	b.IsDone = done
}
