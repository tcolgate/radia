package ghs

type EdgeState int

const (
	EdgeBasic EdgeState = iota
	EdgeBranch
	EdgeRejected
)

type EdgeSet []Edge

type Edge struct {
	Weight Weight
	State  EdgeState
	Core   bool
	Peer   RemoteNode
}

func (EdgeSet) MinEdge() *Edge {
	return nil
}
