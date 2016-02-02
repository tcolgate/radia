package ghs

type EdgeState int

const (
	EdgeBasic EdgeState = iota
	EdgeBranch
	EdgeRejected
)

type Edge struct {
	Weight Weight
	State  EdgeState
	Core   bool

	SenderReciever
}

type Edges []Edge

func (Edges) MinEdge() Edge {
	return Edge{}
}

func NewEdge(f SenderRecieverMaker) (Edge, Edge) {
	c1, c2 := f()
	return Edge{SenderReciever: c1}, Edge{SenderReciever: c2}
}
