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

import "time"

// chanPair is a sender reciever using channels
type chanPair struct {
	send  chan<- MessageMarshaler
	recv  <-chan MessageMarshaler
	delay *time.Duration
}

func (p chanPair) Send(m MessageMarshaler) {
	p.send <- m
}

func (p chanPair) Recieve() (interface{}, error) {
	m := <-p.recv
	if p.delay != nil {
		time.Sleep(*p.delay)
	}

	return m, nil
}

func (p chanPair) Close() {
	close(p.send)
}

// MakeChanPair is an edge sender/reciever built using a
// channel
func MakeChanPair() (SenderReciever, SenderReciever) {
	c1, c2 := make(chan MessageMarshaler), make(chan MessageMarshaler)
	return chanPair{c1, c2, nil}, chanPair{c2, c1, nil}
}

func MakeDelayChanPair(d time.Duration) (SenderReciever, SenderReciever) {
	c1, c2 := make(chan MessageMarshaler), make(chan MessageMarshaler)
	return chanPair{c1, c2, &d}, chanPair{c2, c1, &d}
}
