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

type Sender interface {
	Send(MessageMarshaler)
}

type Reciever interface {
	Recieve() interface{}
}

type Closer interface {
	Close()
}

type SenderReciever interface {
	Sender
	Reciever
	Closer
}

type SenderRecieverMaker func() (SenderReciever, SenderReciever)

type MessageMarshaler interface {
	MarshalMessage() ([]byte, string)
}
