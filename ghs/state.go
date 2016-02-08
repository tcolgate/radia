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
	e.Send(ConnectMessage(n.Level))
}

// Connect - (3) Response to receipt of Connect(L) on edge e
func (n *Node) Connect(e *Edge, level uint32) {
	if n.State == NodeStateSleeping {
		n.procWakeUp()
	}
	if level < n.Level {
		e.State = EdgeStateBranch
		e.Send(InitiateMessage(n.Level, n.Fragment, n.State))
		if n.State == NodeStateFind {
			n.findCount++
		}
	} else if e.State == EdgeStateBasic {
		nm := ConnectMessage(level)
		nm.Edge = e
		n.Queue(nm)
	} else {
		e.Send(InitiateMessage(n.Level+1, e.Weight.FragmentID(), NodeStateFind))
	}
}

// Initiate - (4) Response to receipt of Initiate(L,F,S) on edge e
func (n *Node) Initiate(e *Edge, level uint32, fragment FragmentID, state NodeState) {
	n.Level = level
	n.Fragment = fragment
	n.State = state
	n.inBranch = e
	n.bestEdge = nil
	n.bestWt = WeightInf
	for _, se := range n.Edges {
		if se != e && e.State == EdgeStateBranch {
			se.Send(InitiateMessage(n.Level, n.Fragment, state))
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
			e.Send(TestMessage(n.Level, n.Fragment))
			break
		}
	}
	if !found {
		n.testEdge = nil
		n.procReport()
	}
}

// Test - (5) Response to receipt of Test(L,F) on edge e
func (n *Node) Test(e *Edge, level uint32, fragment FragmentID) {
	if n.State == NodeStateSleeping {
		n.procWakeUp()
	}
	if level > n.Level {
		nm := TestMessage(level, fragment)
		nm.Edge = e
		n.Queue(nm)
	} else {
		if fragment != n.Fragment {
			e.Send(AcceptMessage())
		} else {
			if e.State == EdgeStateBasic {
				e.State = EdgeStateRejected
			}
			if n.testEdge != e {
				e.Send(RejectMessage())
			} else {
				n.procTest()
			}
		}
	}
}

// Accept - (7) Response to receipt of Accept on edge e
func (n *Node) Accept(e *Edge) {
	n.testEdge = nil
	if e.Weight.Less(n.bestWt) {
		n.bestEdge = e
		n.bestWt = e.Weight
	}
	n.procReport()
}

// Reject - (8) Response to receipt of Reject on edge e
func (n *Node) Reject(e *Edge) {
	if e.State == EdgeStateBasic {
		e.State = EdgeStateRejected
	}
	n.procTest()
}

// procReport - (9) procedure report
func (n *Node) procReport() {
	if n.findCount == 0 && n.testEdge == nil {
		n.State = NodeStateFound
		n.inBranch.Send(ReportMessage(n.bestWt))
	}
}

// Report - (10) Response to receipt of Report on edge e
func (n *Node) Report(e *Edge, w Weight) {
	if e != n.inBranch {
		n.findCount--
		if e.Weight.Less(n.bestWt) {
			n.bestWt = w
			n.bestEdge = e
		}
		n.procReport()
	} else if n.State == NodeStateFind {
		nm := ReportMessage(n.bestWt)
		nm.Edge = e
		n.Queue(nm)
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
		n.bestEdge.Send(ChangeRootMessage())
	}
}

// ChangeRoot - (12) Response to receipt of Change-root
func (n *Node) ChangeRoot() {
	n.procChangeRoot()
}
