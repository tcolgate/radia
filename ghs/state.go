package ghs

import "log"

// This file is a translation of the algorithm as described in the
// paper. It's as literal as possible as
// procWakeUp - (2) wakeup node procedure
func (n *Node) procWakeUp() {
	e := n.Edges.MinEdge()
	n.State = NodeStateFound
	n.findCount = 0
	n.Level = 0
	e.Send(Message{
		Func:  (*Node).Connect,
		Level: n.Level,
	})
}

// Connect - (3) Response to receipt of Connect(L) on edge e
func (n *Node) Connect(msg Message) {
	e := msg.Edge
	level := msg.Level

	if n.State == NodeStateSleeping {
		n.procWakeUp()
	}
	if level < n.Level {
		e.State = EdgeStateBranch
		e.Send(Message{
			Func:       (*Node).Initiate,
			Level:      n.Level,
			FragmentID: n.Fragment,
			NodeState:  n.State,
		})
		if n.State == NodeStateFind {
			n.findCount++
		}
	} else if e.State == EdgeStateBasic {
		n.Queue(Message{
			Func:  (*Node).Connect,
			Edge:  e,
			Level: level,
		})
	} else {
		e.Send(Message{
			Func:       (*Node).Initiate,
			Level:      n.Level + 1,
			FragmentID: e.Weight.FragmentID(),
			NodeState:  NodeStateFind,
		})
	}
}

// Initiate - (4) Response to receipt of Initiate(L,F,S) on edge e
func (n *Node) Initiate(msg Message) {
	e := msg.Edge
	level := msg.Level
	fragment := msg.FragmentID
	state := msg.NodeState

	n.Level = level
	n.Fragment = fragment
	n.State = state
	n.inBranch = e
	n.bestEdge = nil
	n.bestWt = WeightInf
	for _, se := range n.Edges {
		if se != e && e.State == EdgeStateBranch {
			se.Send(Message{
				Func:       (*Node).Initiate,
				Level:      n.Level,
				FragmentID: n.Fragment,
				NodeState:  state,
			})
			if n.State == NodeStateFind {
				n.findCount++
			}
		}
		if state == NodeStateFind {
			n.procTest()
		}
	}
}

// ProcTest - (5) procedure test
func (n *Node) procTest() {
	found := false
	for _, e := range n.Edges {
		if e.State == EdgeStateBasic {
			found = true
			n.testEdge = e
			e.Send(Message{
				Func:       (*Node).Test,
				Level:      n.Level,
				FragmentID: n.Fragment,
			})
			break
		}
	}
	if !found {
		n.testEdge = nil
		n.procReport()
	}
}

// Test - (5) Response to receipt of Test(L,F) on edge e
func (n *Node) Test(msg Message) {
	e := msg.Edge
	level := msg.Level
	fragment := msg.FragmentID

	if n.State == NodeStateSleeping {
		n.procWakeUp()
	}
	if level > n.Level {
		n.Queue(Message{
			Func:       (*Node).Test,
			Edge:       e,
			Level:      level,
			FragmentID: fragment,
		})
	} else {
		if fragment != n.Fragment {
			e.Send(Message{
				Func: (*Node).Accept,
			})
		} else {
			if e.State == EdgeStateBasic {
				e.State = EdgeStateRejected
			}
			if n.testEdge != e {
				e.Send(Message{
					Func: (*Node).Reject,
				})
			} else {
				n.procTest()
			}
		}
	}
}

// Accept - (7) Response to receipt of Accept on edge e
func (n *Node) Accept(msg Message) {
	e := msg.Edge

	n.testEdge = nil
	if e.Weight.Less(n.bestWt) {
		n.bestEdge = e
		n.bestWt = e.Weight
	}
	n.procReport()
}

// Reject - (8) Response to receipt of Reject on edge e
func (n *Node) Reject(msg Message) {
	e := msg.Edge

	if e.State == EdgeStateBasic {
		e.State = EdgeStateRejected
	}
	n.procTest()
}

// procReport - (9) procedure report
func (n *Node) procReport() {
	if n.findCount == 0 && n.testEdge == nil {
		n.State = NodeStateFound
		n.inBranch.Send(Message{
			Func:   (*Node).Report,
			Weight: n.bestWt,
		})
	}
}

// Report - (10) Response to receipt of Report on edge e
func (n *Node) Report(msg Message) {
	e := msg.Edge
	w := msg.Weight

	if e != n.inBranch {
		n.findCount--
		if e.Weight.Less(n.bestWt) {
			n.bestWt = w
			n.bestEdge = e
		}
		n.procReport()
	} else if n.State == NodeStateFind {
		n.Queue(Message{
			Func: (*Node).Report,
			Edge: e,
		})
	} else if w.Greater(n.bestWt) {
		n.procChangeRoot()
	} else if w.Equal(n.bestWt) {
		// Halt
		log.Println("Halt")
	}
}

// procChangeRoot - (11) procedure change-root
func (n *Node) procChangeRoot() {
	if n.bestEdge.State == EdgeStateBranch {
		n.bestEdge.Send(Message{
			Func: (*Node).ChangeRoot,
		})
	}
}

// ChangeRoot - (12) Response to receipt of Change-root
func (n *Node) ChangeRoot(msg Message) {
	n.procChangeRoot()
}
