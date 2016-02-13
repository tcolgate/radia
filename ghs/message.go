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

package ghs

import (
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/tcolgate/vonq/graphalg"
)

func init() {
	NodeStateSleeping = NodeState(GHSMessage_Initiate_NodeState_value["sleeping"])
	NodeStateFind = NodeState(GHSMessage_Initiate_NodeState_value["find"])
	NodeStateFound = NodeState(GHSMessage_Initiate_NodeState_value["found"])
}

//go:generate protoc -I $GOPATH/src:. --go_out=.  ghs.proto
type Message struct {
	*GHSMessage
}

func ConnectMessage(level uint32) Message {
	return Message{
		&GHSMessage{
			Type: GHSMessage_CONNECT,
			Connect: &GHSMessage_Connect{
				Level: level,
			},
		},
	}
}

func InitiateMessage(level uint32, fragment FragID, state NodeState) Message {
	wg := graphalg.Weight(fragment)
	return Message{
		&GHSMessage{
			Type: GHSMessage_INITIATE,
			Initiate: &GHSMessage_Initiate{
				Level:     level,
				Fragment:  &wg,
				NodeState: GHSMessage_Initiate_NodeState(state),
			},
		},
	}
}

func TestMessage(level uint32, fragment FragID) Message {
	wg := graphalg.Weight(fragment)
	return Message{
		&GHSMessage{
			Type: GHSMessage_TEST,
			Test: &GHSMessage_Test{
				Level:    level,
				Fragment: &wg,
			},
		},
	}
}

func AcceptMessage() Message {
	return Message{
		GHSMessage: &GHSMessage{
			Type:   GHSMessage_ACCEPT,
			Accept: &GHSMessage_Accept{},
		},
	}
}

func RejectMessage() Message {
	return Message{
		GHSMessage: &GHSMessage{
			Type:   GHSMessage_REJECT,
			Reject: &GHSMessage_Reject{},
		},
	}
}

func ReportMessage(best graphalg.Weight) Message {
	return Message{
		GHSMessage: &GHSMessage{
			Type: GHSMessage_REPORT,
			Report: &GHSMessage_Report{
				Weight: &best,
			},
		},
	}
}

func ChangeRootMessage() Message {
	return Message{
		GHSMessage: &GHSMessage{
			Type:       GHSMessage_CHANGEROOT,
			Changeroot: &GHSMessage_ChangeRoot{},
		},
	}
}

func (s *State) QueueGHS(j int, m Message) {
	b, err := proto.Marshal(m.GHSMessage)
	if err != nil {
		log.Println(err)
	}
	s.Node.Queue(j, b)
}

func (s *State) SendGHS(j int, m Message) {
	b, err := proto.Marshal(m.GHSMessage)
	if err != nil {
		log.Println(err)
	}

	s.Node.Send(j, b)
}

func (s *State) Dispatch(j int, b []byte) {
	m := GHSMessage{}
	proto.Unmarshal(b, &m)

	switch m.Type {
	case GHSMessage_CONNECT:
		s.Connect(j, m.GetConnect().Level)
	case GHSMessage_INITIATE:
		im := m.GetInitiate()
		l := im.Level
		wf := im.Fragment
		ns := NodeState(im.NodeState)
		s.Initiate(j, l, FragmentID(*wf), ns)
	case GHSMessage_TEST:
		im := m.GetTest()
		l := im.Level
		wf := im.Fragment
		s.Test(j, l, FragmentID(*wf))
	case GHSMessage_ACCEPT:
		s.Accept(j)
	case GHSMessage_REJECT:
		s.Reject(j)
	case GHSMessage_REPORT:
		rm := m.GetReport()
		w := rm.GetWeight()
		s.Report(j, *w)
	case GHSMessage_CHANGEROOT:
		s.ChangeRoot()
	default:
		log.Println("unknown message type m.Type")
	}
}
