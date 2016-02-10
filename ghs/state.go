package ghs

import "golang.org/x/net/context"

// This file is a translation of the algorithm as described in the
// paper. It's as literal as possible as

// WakeUp -
func (n *Node) WakeUp(ctx context.Context) {
	n.procWakeUp(ctx)
}

// procWakeUp - (2) wakeup node procedure
func (n *Node) procWakeUp(ctx context.Context) {
	n.Println("Waking")

	e := n.Edges.MinEdge()
	e.State = EdgeStateBranch
	n.Level = 0
	n.State = NodeStateFound
	n.findCount = 0
	e.Send(ctx, ConnectMessage(n.Level))
}

// Connect - (3) Response to receipt of Connect(L) on edge e
func (n *Node) Connect(ctx context.Context, e *Edge, level uint32) {
	if n.State == NodeStateSleeping {
		n.procWakeUp(ctx)
	}
	if level < n.Level {
		e.State = EdgeStateBranch
		e.Send(ctx, InitiateMessage(n.Level, n.Fragment, n.State))
		if n.State == NodeStateFind {
			n.findCount++
		}
	} else {
		if e.State == EdgeStateBasic {
			nm := ConnectMessage(level)
			nm.Edge = e
			n.Queue(nm)
		} else {
			e.Send(ctx, InitiateMessage(n.Level+1, e.Weight.FragmentID(), NodeStateFind))
		}
	}
}

// Initiate - (4) Response to receipt of Initiate(L,F,S) on edge e
func (n *Node) Initiate(ctx context.Context, e *Edge, level uint32, fragment FragmentID, state NodeState) {
	n.Level = level
	n.Fragment = fragment
	n.State = state
	n.inBranch = e
	n.bestEdge = nil
	n.bestWt = WeightInf
	for _, se := range n.Edges {
		if se != e && e.State == EdgeStateBranch {
			se.Send(ctx, InitiateMessage(level, fragment, state))
			if n.State == NodeStateFind {
				n.findCount++
			}
		}
	}
	if state == NodeStateFind {
		n.procTest(ctx)
	}
}

// ProcTest - (5) procedure test
func (n *Node) procTest(ctx context.Context) {
	found := false
	for _, e := range n.Edges {
		if e.State == EdgeStateBasic {
			found = true
			n.testEdge = e
			e.Send(ctx, TestMessage(n.Level, n.Fragment))
			break
		}
	}
	if !found {
		n.testEdge = nil
		n.procReport(ctx)
	}
}

// Test - (6) Response to receipt of Test(L,F) on edge e
func (n *Node) Test(ctx context.Context, e *Edge, level uint32, fragment FragmentID) {
	if n.State == NodeStateSleeping {
		n.procWakeUp(ctx)
	}
	if level > n.Level {
		nm := TestMessage(level, fragment)
		nm.Edge = e
		n.Queue(nm)
	} else {
		if fragment != n.Fragment {
			e.Send(ctx, AcceptMessage())
		} else {
			if e.State == EdgeStateBasic {
				e.State = EdgeStateRejected
			}
			if n.testEdge != e {
				e.Send(ctx, RejectMessage())
			} else {
				n.procTest(ctx)
			}
		}
	}
}

// Accept - (7) Response to receipt of Accept on edge e
func (n *Node) Accept(ctx context.Context, e *Edge) {
	n.testEdge = nil
	if e.Weight.Less(n.bestWt) {
		n.bestEdge = e
		n.bestWt = e.Weight
	}
	n.procReport(ctx)
}

// Reject - (8) Response to receipt of Reject on edge e
func (n *Node) Reject(ctx context.Context, e *Edge) {
	if e.State == EdgeStateBasic {
		e.State = EdgeStateRejected
	}
	n.procTest(ctx)
}

// procReport - (9) procedure report
func (n *Node) procReport(ctx context.Context) {
	if n.findCount == 0 && n.testEdge == nil {
		n.State = NodeStateFound
		n.inBranch.Send(ctx, ReportMessage(n.bestWt))
	}
}

// Report - (10) Response to receipt of Report on edge e
func (n *Node) Report(ctx context.Context, e *Edge, w Weight) {
	if e != n.inBranch {
		n.findCount--
		if w.Less(n.bestWt) {
			n.bestWt = w
			n.bestEdge = e
		}
		n.procReport(ctx)
	} else {
		if n.State == NodeStateFind {
			nm := ReportMessage(w)
			nm.Edge = e
			n.Queue(nm)
		} else {
			if w.Greater(n.bestWt) {
				n.procChangeRoot(ctx)
			} else {
				if w.Equal(n.bestWt) {
					n.Done = true
				}
			}
		}
	}
}

// procChangeRoot - (11) procedure change-root
func (n *Node) procChangeRoot(ctx context.Context) {
	if n.bestEdge.State == EdgeStateBranch {
		n.bestEdge.Send(ctx, ChangeRootMessage())
	} else {
		n.bestEdge.Send(ctx, ConnectMessage(n.Level))
		n.bestEdge.State = EdgeStateBranch
	}
}

// ChangeRoot - (12) Response to receipt of Change-root
func (n *Node) ChangeRoot(ctx context.Context) {
	n.procChangeRoot(ctx)
}
