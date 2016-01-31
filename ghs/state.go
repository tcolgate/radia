package ghs

import "log"

type State func(*Node) State

func Dormant(n *Node) State {
	e := n.Edges.MinEdge()
	log.Println(e)

	return nil
}

// run the state machine
func run(n *Node) {
	for state := Dormant; state != nil; {
		state = state(n)
	}
}
