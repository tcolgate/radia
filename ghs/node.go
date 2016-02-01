package ghs

type NodeID int

type Node struct {
	ID       NodeID
	Edges    EdgeSet
	Level    int
	Fragment FragmentID
}

type RemoteNode struct {
	n1, n2 chan Message
}
