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

//go:generate protoc -I $GOPATH/src:. --go_out=.  ghs.proto
func init() {
	NodeStateSleeping = NodeState(Message_InitiateMsg_NodeState_value["sleeping"])
	NodeStateFind = NodeState(Message_InitiateMsg_NodeState_value["find"])
	NodeStateFound = NodeState(Message_InitiateMsg_NodeState_value["found"])

	typeURL = graphalg.RegisterMessage(Message{})
}

var typeURL = ""

func (m Message) MarshalMessage() ([]byte, string) {
	bs, _ := proto.Marshal(&m)
	return bs, typeURL
}

func ConnectMessage(level uint32) *Message {
	return &Message{
		U: &Message_Connect{
			Connect: &Message_ConnectMsg{
				Level: level,
			},
		},
	}
}

func InitiateMessage(level uint32, fragment FragID, state NodeState) *Message {
	wg := graphalg.Weight(fragment)
	return &Message{
		U: &Message_Initiate{
			Initiate: &Message_InitiateMsg{
				Level:     level,
				Fragment:  &wg,
				NodeState: Message_InitiateMsg_NodeState(state),
			},
		},
	}
}

func TestMessage(level uint32, fragment FragID) *Message {
	wg := graphalg.Weight(fragment)
	return &Message{
		U: &Message_Test{
			Test: &Message_TestMsg{
				Level:    level,
				Fragment: &wg,
			},
		},
	}
}

func AcceptMessage() *Message {
	return &Message{
		U: &Message_Accept{
			Accept: &Message_AcceptMsg{},
		},
	}
}

func RejectMessage() *Message {
	return &Message{
		U: &Message_Reject{
			Reject: &Message_RejectMsg{},
		},
	}
}

func ReportMessage(best graphalg.Weight) *Message {
	return &Message{
		U: &Message_Report{
			Report: &Message_ReportMsg{
				Weight: &best,
			},
		},
	}
}

func ChangeRootMessage() *Message {
	return &Message{
		U: &Message_ChangeRoot{
			ChangeRoot: &Message_ChangeRootMsg{},
		},
	}
}

func HaltMessage() *Message {
	return &Message{
		U: &Message_Halt{
			Halt: &Message_HaltMsg{},
		},
	}
}

func (s *State) QueueGHS(j int, m *Message) {
	s.Queue(j, m)
}

func (s *State) SendGHS(j int, m *Message) {
	s.Send(j, m)
}

func (s *State) Dispatch(j int, i interface{}) {
	m, ok := i.(*Message)
	if !ok {
		log.Fatalf("Non ghs.Message recieved")
	}

	switch m.U.(type) {
	case *Message_Connect:
		s.Connect(j, m.GetConnect().Level)
	case *Message_Initiate:
		im := m.GetInitiate()
		l := im.Level
		wf := im.Fragment
		ns := NodeState(im.NodeState)
		s.Initiate(j, l, FragmentID(*wf), ns)
	case *Message_Test:
		im := m.GetTest()
		l := im.Level
		wf := im.Fragment
		s.Test(j, l, FragmentID(*wf))
	case *Message_Accept:
		s.Accept(j)
	case *Message_Reject:
		s.Reject(j)
	case *Message_Report:
		rm := m.GetReport()
		w := rm.GetWeight()
		s.Report(j, *w)
	case *Message_ChangeRoot:
		s.ChangeRoot()
	case *Message_Halt:
		s.Halt(j)
	default:
		s.Println("unknown message type m.Type")
	}
}
