// Copyright (c) 2016 Tristan Colgate-McFarlane
//
// This file is part of vonq.
//
// vonq is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// vonq is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with vonq.  If not, see <http://www.gnu.org/licenses/>.

package ghs

import "github.com/tcolgate/vonq/graphalg"

// This file is a translation of the algorithm as described in the
// paper. It's as literal as possible as

// NIL is used for identifying unset edge indexes
const NIL = -1

type State struct {
	graphalg.Node
	NodeState
	Fragment FragID
	Level    uint32

	EdgeStates []EdgeState
	bestEdge   int
	testEdge   int
	inBranch   int

	bestWt graphalg.Weight

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

type FragID graphalg.Weight

// FragmentID converts a Weight to a FragmentID. The details of the best
//  edge in a fragment are effectively act as a fragment id.
// In 2) Response to receipt of Connect(... we have
// ...
//    else send Initiate(LN + 1, w(j), Find) on edge j
// ...
// Which clearly uses the edge weight in a Initiate (L, F, S)  message
func FragmentID(w graphalg.Weight) FragID {
	return FragID(w)
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
	s.Send(j, ConnectMessage(s.Level))
}

// Connect - (3) Response to receipt of Connect(L) on edge j
func (s *State) Connect(j int, level uint32) {
	if s.NodeState == NodeStateSleeping {
		s.procWakeUp()
	}
	if level < s.Level {
		s.EdgeStates[j] = EdgeStateBranch
		s.Send(j, InitiateMessage(s.Level, s.Fragment, s.NodeState))
		if s.NodeState == NodeStateFind {
			s.findCount++
		}
	} else {
		if s.EdgeStates[j] == EdgeStateBasic {
			s.Queue(j, ConnectMessage(level))
		} else {
			s.Send(j, InitiateMessage(s.Level+1, FragmentID(s.Edge(j).Weight), NodeStateFind))
		}
	}
}

// Initiate - (4) Response to receipt of Initiate(L,F,S) on edge j
func (s *State) Initiate(j int, level uint32, fragment FragID, state NodeState) {
	s.Level = level
	s.Fragment = fragment
	s.NodeState = state
	s.inBranch = j
	s.bestEdge = NIL
	s.bestWt = graphalg.WeightInf
	for i := range s.EdgeStates {
		if j != i && s.EdgeStates[i] == EdgeStateBranch {
			s.Send(i, InitiateMessage(level, fragment, state))
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
			s.Send(i, TestMessage(s.Level, s.Fragment))
			break
		}
	}
	if !found {
		s.testEdge = NIL
		s.procReport()
	}
}

// Test - (6) Response to receipt of Test(L,F) on edge j
func (s *State) Test(j int, level uint32, fragment FragID) {
	if s.NodeState == NodeStateSleeping {
		s.procWakeUp()
	}
	if level > s.Level {
		s.Queue(j, TestMessage(level, fragment))
	} else {
		if fragment != s.Fragment {
			s.Send(j, AcceptMessage())
		} else {
			if s.EdgeStates[j] == EdgeStateBasic {
				s.EdgeStates[j] = EdgeStateRejected
			}
			if s.testEdge != j {
				s.Send(j, RejectMessage())
			} else {
				s.procTest()
			}
		}
	}
}

// Accept - (7) Response to receipt of Accept on edge j
func (s *State) Accept(j int) {
	s.testEdge = NIL
	if s.Edge(j).Weight.Less(s.bestWt) {
		s.bestEdge = j
		s.bestWt = s.Edge(j).Weight
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
		s.Send(s.inBranch, ReportMessage(s.bestWt))
	}
}

// Report - (10) Response to receipt of Report on edge j
func (s *State) Report(j int, w graphalg.Weight) {
	if j != s.inBranch {
		s.findCount--
		if w.Less(s.bestWt) {
			s.bestWt = w
			s.bestEdge = j
		}
		s.procReport()
	} else {
		if s.NodeState == NodeStateFind {
			s.Queue(j, ReportMessage(w))
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
		s.Send(s.bestEdge, ChangeRootMessage())
	} else {
		s.Send(s.bestEdge, ConnectMessage(s.Level))
		s.EdgeStates[s.bestEdge] = EdgeStateBranch
	}
}

// ChangeRoot - (12) Response to receipt of Change-root
func (s *State) ChangeRoot() {
	s.procChangeRoot()
}

//func (s *State) String() string {
//	return fmt.Sprintf("node(%v)(SN: %v, LN: %v, F: %v, ES: %v, BE: %v, BW: %v, TE: %v, IB: %v, FC: %v)",
//		n.ID, n.State, n.Level, n.Fragment, n.Edges, n.bestEdge, n.bestWt, n.testEdge, n.inBranch, n.findCount)
//}
