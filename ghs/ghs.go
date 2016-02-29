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
	*graphalg.Node

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

// NodeState is the state of the node.
type NodeState Message_InitiateMsg_NodeState

// This is perfect, but these values will be set based
// on what the protoc compiler assigned these values
// but we want nice names for them here
var (
	NodeStateSleeping NodeState
	NodeStateFind     NodeState
	NodeStateFound    NodeState
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
	s.EdgeStates = make([]EdgeState, len(s.Edges()))
	s.Fragment = FragID{Msn: string(s.ID)}

	s.inBranch = NIL
	s.testEdge = NIL
	s.bestEdge = NIL

	j := s.Edges().MinEdge()
	s.EdgeStates[j] = EdgeStateBranch
	s.Level = 0
	s.NodeState = NodeStateFound
	s.findCount = 0
	s.SendGHS(j, ConnectMessage(s.Level))
}

// Connect - (3) Response to receipt of Connect(L) on edge j
func (s *State) Connect(j int, level uint32) {
	if s.NodeState == NodeStateSleeping {
		s.procWakeUp()
	}
	if level < s.Level {
		s.EdgeStates[j] = EdgeStateBranch
		s.SendGHS(j, InitiateMessage(s.Level, s.Fragment, s.NodeState))
		if s.NodeState == NodeStateFind {
			s.findCount++
		}
	} else {
		if s.EdgeStates[j] == EdgeStateBasic {
			s.QueueGHS(j, ConnectMessage(level))
		} else {
			s.SendGHS(j, InitiateMessage(s.Level+1, FragmentID(s.Edge(j).Weight), NodeStateFind))
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
			s.SendGHS(i, InitiateMessage(level, fragment, state))
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
			s.SendGHS(i, TestMessage(s.Level, s.Fragment))
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
		s.QueueGHS(j, TestMessage(level, fragment))
	} else {
		if fragment != s.Fragment {
			s.SendGHS(j, AcceptMessage())
		} else {
			if s.EdgeStates[j] == EdgeStateBasic {
				s.EdgeStates[j] = EdgeStateRejected
			}
			if s.testEdge != j {
				s.SendGHS(j, RejectMessage())
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
		s.SendGHS(s.inBranch, ReportMessage(s.bestWt))
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
			s.QueueGHS(j, ReportMessage(w))
		} else {
			if w.Greater(s.bestWt) {
				s.procChangeRoot()
			} else {
				if w.Equal(s.bestWt) {
					// We deviate from the original algorithm here
					// as originall only core nodes halt. We'll
					// broadcast a halt to all nodes.
					s.procHalt()
				}
			}
		}
	}
}

// procChangeRoot - (11) procedure change-root
func (s *State) procChangeRoot() {
	if s.EdgeStates[s.bestEdge] == EdgeStateBranch {
		s.SendGHS(s.bestEdge, ChangeRootMessage())
	} else {
		s.SendGHS(s.bestEdge, ConnectMessage(s.Level))
		s.EdgeStates[s.bestEdge] = EdgeStateBranch
	}
}

// ChangeRoot - (12) Response to receipt of Change-root
func (s *State) ChangeRoot() {
	s.procChangeRoot()
}

// ProcHalt - Tell all nodes we have finished
func (s *State) procHalt() {
	for i := range s.EdgeStates {
		if i != s.inBranch && s.EdgeStates[i] == EdgeStateBranch &&
			FragmentID(s.Edge(i).Weight) != s.Fragment {
			s.SendGHS(i, HaltMessage())
		}
	}
	s.SetDone(true)
	return
}

// Halt - Response to receipt of Halt
func (s *State) Halt(j int) {
	s.inBranch = j
	s.procHalt()
}
