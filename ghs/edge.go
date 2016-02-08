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

func (e *Edge) send(m Message) {
	e.SenderReciever.Send(m)
}

func (e Edge) SendConnect(level int) {
	e.send(Message{
		Type:  MessageConnect,
		Level: level,
	})
}

func (e Edge) SendInitiate(level int, fragment FragmentID, state NodeState) {
	e.send(Message{
		Type:       MessageInitiate,
		Level:      level,
		FragmentID: fragment,
		NodeState:  state,
	})
}

func (e Edge) SendTest(level int, fragment FragmentID) {
	e.send(Message{
		Type:       MessageTest,
		Level:      level,
		FragmentID: fragment,
	})
}

func (e Edge) SendAccept() {
	e.send(Message{
		Type: MessageAccept,
	})
}

func (e Edge) SendReject() {
	e.send(Message{
		Type: MessageReject,
	})
}

func (e Edge) SendReport(best Weight) {
	e.send(Message{
		Type:   MessageReport,
		Weight: best,
	})
}

func (e Edge) SendChangeRoot() {
	e.send(Message{
		Type: MessageChangeRoot,
	})
}
