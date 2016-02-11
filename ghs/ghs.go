package ghs

import (
	"fmt"

	"github.com/tcolgate/vonq/graphalg"
)

// This file is a translation of the algorithm as described in the
// paper. It's as literal as possible as

const NIL = -1

type State struct {
	graphalg.Node
	NodeState
	Fragment FragmentID
	Level    uint32

	EdgeStates []EdgeState
	bestEdge   int
	testEdge   int
	inBranch   int

	bestWt Weight

	findCount int
}

//go:generate stringer -type=NodeState
type NodeState int

const (
	NodeStateSleeping NodeState = iota
	NodeStateFind
	NodeStateFound
)

//go:generate stringer -type=EdgeState
type EdgeState int

const (
	EdgeStateBasic EdgeState = iota
	EdgeStateBranch
	EdgeStateRejected
)

// FragmentID converts a Weight to a FragmentID. The details of the best
//  edge in a fragment are effectively act as a fragment id.
// In 2) Response to receipt of Connect(... we have
// ...
//    else send Initiate(LN + 1, w(j), Find) on edge j
// ...
// Which clearly uses the edge weight in a Initiate (L, F, S)  message
func FragmentID(w graphalg.Weight) FragmentID {
	return FragmentID(w)
}

// WakeUp -
func (s *State) WakeUp() {
	s.procWakeUp()
}

// procWakeUp - (2) wakeup node procedure
func (s *State) procWakeUp() {
	s.Println("Waking")

	j := s.Edges.MinEdge()
	s.EdgeStates[j] = EdgeStateBranch
	s.Level = 0
	s.NodeState = NodeStateFound
	s.findCount = 0
	s.Edges(j).Send(ConnectMessage(s.Level))
}

// Connect - (3) Response to receipt of Connect(L) on edge j
func (s *State) Connect(j int, level uint32) {
	if s.NodeState == NodeStateSleeping {
		s.procWakeUp()
	}
	if level < s.Level {
		s.EdgesStates[j] = EdgeStateBranch
		s.Edges(j).Send(InitiateMessage(s.Level, s.Fragment, s.NodeState))
		if s.NodeState == NodeStateFind {
			s.findCount++
		}
	} else {
		if s.EdgeStates[j] == EdgeStateBasic {
			nm := ConnectMessage(level)
			nm.Edge = j
			s.Queue(nm)
		} else {
			s.Edges(j).Send(InitiateMessage(s.Level+1, s.Edges(j).Weight.FragmentID(), NodeStateFind))
		}
	}
}

// Initiate - (4) Response to receipt of Initiate(L,F,S) on edge j
func (s *State) Initiate(j int, level uint32, fragment FragmentID, state NodeState) {
	s.Level = level
	s.Fragment = fragment
	s.NodeState = state
	s.inBranch = j
	s.bestEdge = NIL
	s.bestWt = WeightInf
	for i := range s.EdgeStates {
		if j != i && s.EdgeStates[i] == EdgeStateBranch {
			s.Edges(j).Send(InitiateMessage(level, fragment, state))
			if s.NodeState == NodeStateFind {
				s.findCount++
			}
		}
	}
	if state == NodeStateFind {
		s.procTest()
	}
}

// ProcTest - (5) procedure test
func (s *State) procTest() {
	found := false
	for i := range s.EdgeStates {
		if s.EdgeStates[i] == EdgeStateBasic {
			found = true
			s.testEdge = i
			s.Edges(i).Send(TestMessage(s.Level, s.Fragment))
			break
		}
	}
	if !found {
		s.testEdge = NIL
		s.procReport()
	}
}

// Test - (6) Response to receipt of Test(L,F) on edge j
func (s *State) Test(j int, level uint32, fragment FragmentID) {
	if s.NodeState == NodeStateSleeping {
		s.procWakeUp()
	}
	if level > s.Level {
		nm := TestMessage(level, fragment)
		nm.Edge = j
		s.Queue(nm)
	} else {
		if fragment != s.Fragment {
			s.Edges(j).Send(AcceptMessage())
		} else {
			if s.EdgeStates[j] == EdgeStateBasic {
				s.EdgeStates[j] = EdgeStateRejected
			}
			if s.testEdge != j {
				s.Edges(j).Send(RejectMessage())
			} else {
				s.procTest()
			}
		}
	}
}

// Accept - (7) Response to receipt of Accept on edge j
func (s *State) Accept(j int) {
	s.testEdge = NIL
	if s.Edges(j).Weight.Less(s.bestWt) {
		s.bestEdge = j
		s.bestWt = s.Edges(j).Weight
	}
	s.procReport()
}

// Reject - (8) Response to receipt of Reject on edge j
func (s *State) Reject(j int) {
	if s.EdgeStates[j] == EdgeStateBasic {
		s.EdgeStates[j] = EdgeStateRejected
	}
	s.procTest()
}

// procReport - (9) procedure report
func (s *State) procReport() {
	if s.findCount == 0 && s.testEdge == NIL {
		s.NodeState = NodeStateFound
		s.Edges(s.inBranch).Send(ReportMessage(s.bestWt))
	}
}

// Report - (10) Response to receipt of Report on edge j
func (s *State) Report(j int, w Weight) {
	if j != s.inBranch {
		s.findCount--
		if w.Less(s.bestWt) {
			s.bestWt = w
			s.bestEdge = j
		}
		s.procReport()
	} else {
		if s.NodeState == NodeStateFind {
			nm := ReportMessage(w)
			nm.Edge = j
			s.Queue(nm)
		} else {
			if w.Greater(s.bestWt) {
				s.procChangeRoot()
			} else {
				if w.Equal(s.bestWt) {
					s.Done = true
				}
			}
		}
	}
}

// procChangeRoot - (11) procedure change-root
func (s *State) procChangeRoot() {
	if s.EdgeStates[s.bestEdge] == EdgeStateBranch {
		s.Edges(j).Send(ChangeRootMessage())
	} else {
		s.Edges(j).Send(ConnectMessage(s.Level))
		s.EdgeStates[s.bestEdge] = EdgeStateBranch
	}
}

// ChangeRoot - (12) Response to receipt of Change-root
func (s *State) ChangeRoot() {
	s.procChangeRoot()
}

func (s *State) String() string {
	return fmt.Sprintf("node(%v)(SN: %v, LN: %v, F: %v, ES: %v, BE: %v, BW: %v, TE: %v, IB: %v, FC: %v)",
		n.ID, n.State, n.Level, n.Fragment, n.Edges, n.bestEdge, n.bestWt, n.testEdge, n.inBranch, n.findCount)
}
