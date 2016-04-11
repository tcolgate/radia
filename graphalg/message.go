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

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	google_protobuf "github.com/golang/protobuf/ptypes/any"
)

type Sender interface {
	Send(MessageMarshaler)
}

type Reciever interface {
	Recieve() (interface{}, error)
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

var (
	typeRegistry = make(map[string]reflect.Type)
)

func RegisterMessage(i interface{}) string {
	t := reflect.TypeOf(i)
	ti := reflect.New(t).Interface()
	url := fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
	if _, ok := ti.(proto.Message); !ok {
		panic("type does " + url + "not support proto.Message interface")
	}
	typeRegistry[url] = t
	return url
}

type MessageMarshaler interface {
	MarshalMessage() ([]byte, string)
}

func unmarshalAny(any *google_protobuf.Any) (interface{}, error) {
	class := any.TypeUrl
	bytes := any.Value

	t, ok := typeRegistry[class]
	if !ok {
		panic("No type registered for " + class)
	}

	instance := reflect.New(t).Interface()
	err := proto.Unmarshal(bytes, instance.(proto.Message))

	return instance, err
}
