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

package graphalg

import (
	"math"
	"strings"
)

// Weight represents the weight of an edge. The edge includes
// two NodeIDs to break ties in cases where all edges must have
// unqiue weights (e.g. MSTs)
//type Weight struct {
//	Cost float64
//	Lsn  NodeID // Least Signigicant NodeID
//	Msn  NodeID // Most Signigicant NodeID
//}

// WeightInf is an edge with infinite weight
var WeightInf = Weight{Cost: math.Inf(1)}

// Less - compare two edge weights.
func (w1 Weight) Less(w2 Weight) bool {
	switch {
	case w1.Cost < w2.Cost:
		return true
	case w1.Cost == w2.Cost &&
		strings.Compare(string(w1.Msn), string(w2.Msn)) < 0:
		return true
	case w1.Cost == w2.Cost &&
		strings.Compare(string(w1.Msn), string(w2.Msn)) == 0 &&
		strings.Compare(string(w1.Lsn), string(w2.Lsn)) < 0:
		return true
	default:
		return false
	}
}

// Less - compare two edge weights.
func (w1 Weight) Greater(w2 Weight) bool {
	switch {
	case w1.Cost > w2.Cost:
		return true
	case w1.Cost == w2.Cost &&
		strings.Compare(string(w1.Msn), string(w2.Msn)) > 0:
		return true
	case w1.Cost == w2.Cost &&
		strings.Compare(string(w1.Msn), string(w2.Msn)) == 0 &&
		strings.Compare(string(w1.Lsn), string(w2.Lsn)) > 0:
		return true
	default:
		return false
	}
}

// Less - compare two edge weights.
func (w1 Weight) Equal(w2 Weight) bool {
	return w1.Cost == w2.Cost &&
		strings.Compare(string(w1.Lsn), string(w2.Lsn)) == 0 &&
		strings.Compare(string(w1.Msn), string(w2.Msn)) == 0
}
