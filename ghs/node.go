package ghs

type Node struct {
	Edges    EdgeSet
	Level    int
	Fragment FragmentID
}

type RemoteNode chan Message
