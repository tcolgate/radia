package ghs

import "log"
import pb "github.com/tcolgate/vonq/ghs/proto"

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

type Message struct {
	pb.GHSMessage
	Edge *Edge
}

type SenderRecieverMaker func() (SenderReciever, SenderReciever)

func pbWeightToWeight(*pb.GHSMessage_Weight) Weight {
	return Weight{}
}

func pbNodeStateToNodeState(pb.GHSMessage_Initiate_NodeState) NodeState {
	return NodeStateFind
}

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
	return Message{}
}

func InitiateMessage(level uint32, fragment FragmentID, state NodeState) Message {
	/*
		e.send(Message{
			Type:       MessageInitiate,
			Level:      level,
			FragmentID: fragment,
			NodeState:  state,
		})
	*/
	return Message{}
}

func TestMessage(level uint32, fragment FragmentID) Message {
	/*
		e.send(Message{
			Type:       MessageTest,
			Level:      level,
			FragmentID: fragment,
		})
	*/
	return Message{}
}

func AcceptMessage() Message {
	/*
		e.send(Message{
			Type: MessageAccept,
		})
	*/
	return Message{}
}

func RejectMessage() Message {
	/*
		e.send(Message{
			Type: MessageReject,
		})
	*/
	return Message{}
}

func ReportMessage(best Weight) Message {
	/*
		e.send(Message{
			Type:   MessageReport,
			Weight: best,
		})
	*/
	return Message{}
}

func ChangeRootMessage() Message {
	/*
		e.send(Message{
			Type: MessageChangeRoot,
		})
	*/
	return Message{}
}
