package ghs

import (
	"log"
	"sort"
)

type NodeID string

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
	Level    uint32
	Fragment FragmentID
	Done     bool
	OnDone   func()

	msgQueue []Message

	bestEdge  *Edge
	bestWt    Weight
	testEdge  *Edge
	inBranch  *Edge
	findCount int
}

func Join(n1 *Node, n2 *Node, w float64, f SenderRecieverMaker) {
	n1.Join(n2, w, f)
}

func (n1 *Node) Join(n2 *Node, w float64, f SenderRecieverMaker) {
	ids := []string{
		string(n1.ID),
		string(n2.ID),
	}
	sort.Strings(ids)

	e1, e2 := NewEdge(f)
	e1.Weight.float64 = w
	e2.Weight.float64 = w
	e1.Weight.Lsn = NodeID(ids[0])
	e2.Weight.Lsn = NodeID(ids[0])
	e1.Weight.Msn = NodeID(ids[1])
	e2.Weight.Msn = NodeID(ids[1])

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
			log.Printf("node(%v).Edge(%b): Listening", n.ID, *e)
			for {
				ms <- e.Recieve()
			}
		}(e)
	}

	for nm := range ms {
		delayed := n.msgQueue
		n.msgQueue = []Message{}
		log.Printf("node(%v) Dispathing: %v\n", n.ID, nm)
		nm.dispatch(n)

		for _, om := range delayed {
			log.Printf("node(%v) Replaying: %v\n", n.ID, om)
			om.dispatch(n)
			if n.Done {
				if n.OnDone != nil {
					n.OnDone()
				}
				return
			}
		}
	}
}
