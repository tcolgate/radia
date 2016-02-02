package ghs

import "log"

type State func(*Node) State

func Dormant(n *Node) State {
	e := n.Edges.MinEdge()
	e.State = EdgeBranch
	msg := e.Recieve()

	log.Println(msg)
	return nil
}

// run the state machine
func run(n *Node) {
	for state := Dormant; state != nil; {
		state = state(n)
	}
}
