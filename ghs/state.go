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
	e.SendConnect(n.Level)
}

// Connect - (3) Response to receipt of Connect(L) on edge e
func (n *Node) Connect(e *Edge, level int) {
	if n.State == NodeStateSleeping {
		n.procWakeUp()
	}
	if level < n.Level {
		e.State = EdgeStateBranch
		e.SendInitiate(n.Level, n.Fragment, n.State)
		if n.State == NodeStateFind {
			n.findCount++
		}
	} else if e.State == EdgeStateBasic {
		n.Queue(Message{
			Type:  MessageConnect,
			Edge:  e,
			Level: level,
		})
	} else {
		e.SendInitiate(n.Level+1, e.Weight.FragmentID(), NodeStateFind)
	}
}

// Initiate - (4) Response to receipt of Initiate(L,F,S) on edge e
func (n *Node) Initiate(e *Edge, level int, fragment FragmentID, state NodeState) {
	n.Level = level
	n.Fragment = fragment
	n.State = state
	n.inBranch = e
	n.bestEdge = nil
	n.bestWt = WeightInf
	for _, se := range n.Edges {
		if se != e && e.State == EdgeStateBranch {
			se.SendInitiate(n.Level, n.Fragment, state)
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
	n.Edges.SortByMinEdge()
	found := false
	for _, e := range n.Edges {
		if e.State == EdgeStateBasic {
			found = true
			n.testEdge = e
			e.SendTest(n.Level, n.Fragment)
			break
		}
	}
	if !found {
		n.testEdge = nil
		n.procReport()
	}
}

// Test - (5) Response to receipt of Test(L,F) on edge e
func (n *Node) Test(e *Edge, level int, fragment FragmentID) {
	if n.State == NodeStateSleeping {
		n.procWakeUp()
	}
	if level > n.Level {
		n.Queue(Message{
			Type:       MessageTest,
			Edge:       e,
			Level:      level,
			FragmentID: fragment,
		})
	} else {
		if fragment != n.Fragment {
			e.SendAccept()
		} else {
			if e.State == EdgeStateBasic {
				e.State = EdgeStateRejected
			}
			if n.testEdge != e {
				e.SendReject()
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
		n.inBranch.SendReport(n.bestWt)
	}
}

// Report - (10)
func (n *Node) Report(e *Edge, w Weight) {
	if e != n.inBranch {
		n.findCount--
		if e.Weight.Less(n.bestWt) {
			n.bestWt = w
			n.bestEdge = e
		}
		n.procReport()
	} else if n.State == NodeStateFind {
		n.Queue(Message{
			Type: MessageReport,
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
		n.bestEdge.SendChangeRoot()
	}
}

// ChangeRoot - (12) Response to receipt of Change-root
func (n *Node) ChangeRoot(e *Edge) {
	n.procChangeRoot()
}
