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
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	pb "github.com/tcolgate/vonq/ghs/proto"
	"github.com/tcolgate/vonq/graphalg"
)

type Message struct {
	*graphalg.Edge
	pb.GHSMessage
}

func pbWeightToWeight(w *pb.GHSMessage_Weight) Weight {
	return Weight{
		float64: w.GetWeight(),
		Lsn:     NodeID(w.GetLsnID()),
		Msn:     NodeID(w.GetMsnID()),
	}
}

func pbNodeStateToNodeState(n pb.GHSMessage_Initiate_NodeState) NodeState {
	switch n.String() {
	case "sleeping":
		return NodeStateSleeping
	case "find":
		return NodeStateFind
	case "found":
		return NodeStateFound
	}
	panic("eah?")
}

func protoWeight(w Weight) *pb.GHSMessage_Weight {
	return &pb.GHSMessage_Weight{
		Weight: proto.Float64(w.float64),
		LsnID:  proto.String(string(w.Lsn)),
		MsnID:  proto.String(string(w.Msn)),
	}
}

func protoFragmentID(f FragmentID) *pb.GHSMessage_Weight {
	return &pb.GHSMessage_Weight{
		Weight: proto.Float64(f.float64),
		LsnID:  proto.String(string(f.Lsn)),
		MsnID:  proto.String(string(f.Msn)),
	}
}

func protoNodeState(n NodeState) *pb.GHSMessage_Initiate_NodeState {
	switch n {
	case NodeStateSleeping:
		return pb.GHSMessage_Initiate_sleeping.Enum()
	case NodeStateFind:
		return pb.GHSMessage_Initiate_find.Enum()
	case NodeStateFound:
		return pb.GHSMessage_Initiate_found.Enum()
	}
	return nil
}

func (m Message) String() string {
	switch m.GetType() {
	case pb.GHSMessage_CONNECT:
		return fmt.Sprintf("(CONNECT %s)", m.Connect)
	case pb.GHSMessage_INITIATE:
		return fmt.Sprintf("(INITIATE %s)", m.Initiate)
	case pb.GHSMessage_TEST:
		return fmt.Sprintf("(TEST %s)", m.Test)
	case pb.GHSMessage_ACCEPT:
		return fmt.Sprintf("(ACCEPT %s)", m.Accept)
	case pb.GHSMessage_REJECT:
		return fmt.Sprintf("(REJECT %s)", m.Reject)
	case pb.GHSMessage_REPORT:
		return fmt.Sprintf("(REPORT %s)", m.Report)
	case pb.GHSMessage_CHANGEROOT:
		return fmt.Sprintf("(CHANGEROOT %s)", m.Changeroot)
	default:
		return fmt.Sprintf("(unknown message)")
	}
}

func ConnectMessage(level uint32) Message {
	return Message{
		GHSMessage: pb.GHSMessage{
			Type: pb.GHSMessage_CONNECT.Enum(),
			Connect: &pb.GHSMessage_Connect{
				Level: proto.Uint32(level),
			},
		},
	}
}

func InitiateMessage(level uint32, fragment FragmentID, state NodeState) Message {
	return Message{
		GHSMessage: pb.GHSMessage{
			Type: pb.GHSMessage_INITIATE.Enum(),
			Initiate: &pb.GHSMessage_Initiate{
				Level:     proto.Uint32(level),
				Fragment:  protoFragmentID(fragment),
				NodeState: protoNodeState(state),
			},
		},
	}
}

func TestMessage(level uint32, fragment FragmentID) Message {
	return Message{
		GHSMessage: pb.GHSMessage{
			Type: pb.GHSMessage_TEST.Enum(),
			Test: &pb.GHSMessage_Test{
				Level:    proto.Uint32(level),
				Fragment: protoFragmentID(fragment),
			},
		},
	}
}

func AcceptMessage() Message {
	return Message{
		GHSMessage: pb.GHSMessage{
			Type:   pb.GHSMessage_ACCEPT.Enum(),
			Accept: &pb.GHSMessage_Accept{},
		},
	}
}

func RejectMessage() Message {
	return Message{
		GHSMessage: pb.GHSMessage{
			Type:   pb.GHSMessage_REJECT.Enum(),
			Reject: &pb.GHSMessage_Reject{},
		},
	}
}

func ReportMessage(best Weight) Message {
	return Message{
		GHSMessage: pb.GHSMessage{
			Type: pb.GHSMessage_REPORT.Enum(),
			Report: &pb.GHSMessage_Report{
				Weight: protoWeight(best),
			},
		},
	}
}

func ChangeRootMessage() Message {
	return Message{
		GHSMessage: pb.GHSMessage{
			Type:       pb.GHSMessage_CHANGEROOT.Enum(),
			Changeroot: &pb.GHSMessage_ChangeRoot{},
		},
	}
}

func (s *State) Dispatch(m Message) {
	switch m.GetType() {
	case pb.GHSMessage_CONNECT:
		s.Connect(m.Edge, m.GetConnect().GetLevel())
	case pb.GHSMessage_INITIATE:
		im := m.GetInitiate()
		l := im.GetLevel()
		wf := pbWeightToWeight(im.GetFragment())
		ns := pbNodeStateToNodeState(im.GetNodeState())
		s.Initiate(m.Edge, l, wf.FragmentID(), ns)
	case pb.GHSMessage_TEST:
		im := m.GetTest()
		l := im.GetLevel()
		wf := pbWeightToWeight(im.GetFragment())
		s.Test(m.Edge, l, wf.FragmentID())
	case pb.GHSMessage_ACCEPT:
		s.Accept(m.Edge)
	case pb.GHSMessage_REJECT:
		s.Reject(m.Edge)
	case pb.GHSMessage_REPORT:
		rm := m.GetReport()
		w := pbWeightToWeight(rm.GetWeight())
		s.Report(m.Edge, w)
	case pb.GHSMessage_CHANGEROOT:
		s.ChangeRoot()
	default:
		log.Println("unknown message type m.Type")
	}
}
