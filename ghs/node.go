package ghs

import (
	"fmt"
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

func (n *Node) String() string {
	return fmt.Sprintf("node(%v)(SN: %v, LN: %v, F: %v, ES: %v, BE: %v, BW: %v, TE: %v, IB: %v, FC: %v)",
		n.ID, n.State, n.Level, n.Fragment, n.Edges, n.bestEdge, n.bestWt, n.testEdge, n.inBranch, n.findCount)
}

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

	*log.Logger
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

	e1.local, e1.remote = n1, n2 // mostly for debugging
	e2.local, e2.remote = n2, n1

	n1.Edges = append(n1.Edges, e1)
	n2.Edges = append(n2.Edges, e2)
}

// Queue - add a GHS message to the internal queue
func (n *Node) Queue(msg Message) {
	n.Printf("Queueing  %v\n", msg)
	n.msgQueue = append(n.msgQueue, msg)
}

func (n *Node) Run() {
	ms := make(chan Message)
	defer close(ms)

	n.Edges.SortByMinEdge()
	defer func() {
		if n.OnDone != nil {
			n.OnDone()
		}
	}()

	for _, e := range n.Edges {
		go func(e *Edge) {
			n.Printf(".Edge(%b): Listening", *e)
			for {
				ms <- e.Recieve()
			}
		}(e)
	}

	for nm := range ms {
		delayed := n.msgQueue
		n.msgQueue = []Message{}
		n.Printf("before %+v\n", n)
		n.Printf("Do %+v\n", nm)
		nm.dispatch(n)
		n.Printf("after %+v\n", n)
		for _, om := range delayed {
			n.Printf("Redo %+v\n", om)
			om.dispatch(n)
			n.Printf("%+v\n", n)
		}
		if n.Done {
			return
		}
	}
}
