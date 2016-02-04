package ghs

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

func (Edges) MinEdge() *Edge {
	return &Edge{}
}

func NewEdge(f SenderRecieverMaker) (*Edge, *Edge) {
	c1, c2 := f()
	return &Edge{SenderReciever: c1}, &Edge{SenderReciever: c2}
}

func (e Edge) SendConnect(level int) {
}

func (e Edge) SendInitiate(level int, fragment FragmentID, state NodeState) {
}

func (e Edge) SendTest(level int, fragment FragmentID) {
}

func (e Edge) SendAccept() {
}

func (e Edge) SendReject() {
}

func (e Edge) SendReport(best Weight) {
}
