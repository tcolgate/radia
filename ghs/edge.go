package ghs

type EdgeState int

const (
	Unknown EdgeState = iota
	Basic
	Rejected
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
