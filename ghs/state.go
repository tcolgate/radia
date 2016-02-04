package ghs

type State func(*Node) State

func Dormant(n *Node) State {
	return nil
}

// run the state machine
func run(n *Node) {
	for state := Dormant; state != nil; {
		state = state(n)
	}
}
