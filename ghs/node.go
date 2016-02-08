package ghs

import "log"

type NodeID int

//go:generate stringer -type=NodeState
type NodeState int

const (
	NodeStateSleeping NodeState = iota
	NodeStateFind
	NodeStateFound
)

type Node struct {
	ID       NodeID
	State    NodeState
	Edges    Edges
	Level    int
	Fragment FragmentID
	Done     bool

	msgQueue []Message

	bestEdge  *Edge
	bestWt    Weight
	testEdge  *Edge
	inBranch  *Edge
	findCount int
}

func Join(n1 *Node, n2 *Node, f SenderRecieverMaker) {
	e1, e2 := NewEdge(f)
	n1.Edges = append(n1.Edges, e1)
	n2.Edges = append(n2.Edges, e2)
}

func (n1 *Node) Join(n2 *Node, f SenderRecieverMaker) {
	e1, e2 := NewEdge(f)
	n1.Edges = append(n1.Edges, e1)
	n2.Edges = append(n2.Edges, e2)
}

// Queue - add a GHS message to the internal queue
func (n *Node) Queue(msg Message) {
	n.msgQueue = append(n.msgQueue, msg)
}

func (n *Node) Run() {
	ms := make(chan Message)
	n.Edges.SortByMinEdge()

	for _, e := range n.Edges {
		go func(e *Edge) {
			for {
				ms <- e.Recieve()
			}
		}(e)
	}

	for nm := range ms {
		delayed := n.msgQueue
		n.msgQueue = []Message{}
		n.dispatch(nm)

		for _, om := range delayed {
			n.dispatch(om)
			if n.Done {
				return
			}
		}
	}
}

func (n *Node) dispatch(m Message) {
	switch m.Type {
	case MessageConnect:
		n.Connect(m)
	case MessageInitiate:
		n.Initiate(m)
	case MessageTest:
		n.Test(m)
	case MessageAccept:
		n.Accept(m)
	case MessageReject:
		n.Reject(m)
	case MessageReport:
		n.Report(m)
	case MessageChangeRoot:
		n.ChangeRoot(m)
	default:
		log.Println("unknown message type m.Type")
	}
}
