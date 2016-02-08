package ghs

import "sort"

//go:generate stringer -type=EdgeState
type EdgeState int

const (
	EdgeStateBasic EdgeState = iota
	EdgeStateBranch
	EdgeStateRejected
)

type Edge struct {
	Weight Weight
	State  EdgeState

	SenderReciever
}

type Edges []*Edge

func (e Edges) MinEdge() *Edge {
	if len(e) == 0 {
		return nil
	}
	e.SortByMinEdge()
	return e[0]
}

func (e Edges) SortByMinEdge() {
	sort.Sort(ByMinEdge(e))
}

// ByMinEdge - implements sort by minimum edge
type ByMinEdge []*Edge

func (e ByMinEdge) Len() int {
	return len(e)
}

func (e ByMinEdge) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
func (e ByMinEdge) Less(i, j int) bool {
	return e[i].Weight.Less(e[j].Weight)
}

func NewEdge(f SenderRecieverMaker) (*Edge, *Edge) {
	c1, c2 := f()
	return &Edge{SenderReciever: c1}, &Edge{SenderReciever: c2}
}

func (e *Edge) Recieve() Message {
	m := e.SenderReciever.Recieve()
	m.Edge = e
	return m
}

func (e *Edge) Send(m Message) {
	e.SenderReciever.Send(m)
}

// chanPair is a sender reciever using channels
type chanPair struct {
	send chan<- Message
	recv <-chan Message
}

func (p chanPair) Send(m Message) {
	p.send <- m
}

func (p chanPair) Recieve() Message {
	return <-p.recv
}

func (p chanPair) Close() {
	close(p.send)
}

func MakeChanPair() (SenderReciever, SenderReciever) {
	c1, c2 := make(chan Message), make(chan Message)
	return chanPair{c1, c2}, chanPair{c2, c1}
}
