package ghs

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	pb "github.com/tcolgate/vonq/ghs/proto"
)

type Sender interface {
	Send(Message)
}

type Reciever interface {
	Recieve() Message
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

type Message struct {
	pb.GHSMessage
	Edge *Edge
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

// This can probably be shift around again
// Not sure if GHS messages need to know so much
// about the protocol. Maybe it's OK
func (m Message) dispatch(n *Node) {
	switch m.GetType() {
	case pb.GHSMessage_CONNECT:
		n.Connect(m.Edge, m.GetConnect().GetLevel())
	case pb.GHSMessage_INITIATE:
		im := m.GetInitiate()
		l := im.GetLevel()
		wf := pbWeightToWeight(im.GetFragment())
		s := pbNodeStateToNodeState(im.GetNodeState())
		n.Initiate(m.Edge, l, wf.FragmentID(), s)
	case pb.GHSMessage_TEST:
		im := m.GetTest()
		l := im.GetLevel()
		wf := pbWeightToWeight(im.GetFragment())
		n.Test(m.Edge, l, wf.FragmentID())
	case pb.GHSMessage_ACCEPT:
		n.Accept(m.Edge)
	case pb.GHSMessage_REJECT:
		n.Reject(m.Edge)
	case pb.GHSMessage_REPORT:
		rm := m.GetReport()
		w := pbWeightToWeight(rm.GetWeight())
		n.Report(m.Edge, w)
	case pb.GHSMessage_CHANGEROOT:
		n.ChangeRoot()
	default:
		log.Println("unknown message type m.Type")
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
