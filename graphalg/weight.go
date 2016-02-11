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

type Weight struct {
	float64
	Lsn NodeID // Least Signigicant NodeID
	Msn NodeID // Most Signigicant NodeID
}

var WeightInf = Weight{float64: math.Inf(1)}

// Less - compare two edge weights.
func (w1 Weight) Less(w2 Weight) bool {
	switch {
	case w1.float64 < w2.float64:
		return true
	case w1.float64 == w2.float64 &&
		strings.Compare(string(w1.Msn), string(w2.Msn)) < 0:
		return true
	case w1.float64 == w2.float64 &&
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
	case w1.float64 > w2.float64:
		return true
	case w1.float64 == w2.float64 &&
		strings.Compare(string(w1.Msn), string(w2.Msn)) > 0:
		return true
	case w1.float64 == w2.float64 &&
		strings.Compare(string(w1.Msn), string(w2.Msn)) == 0 &&
		strings.Compare(string(w1.Lsn), string(w2.Lsn)) > 0:
		return true
	default:
		return false
	}
}

// Less - compare two edge weights.
func (w1 Weight) Equal(w2 Weight) bool {
	return w1.float64 == w2.float64 &&
		strings.Compare(string(w1.Lsn), string(w2.Lsn)) == 0 &&
		strings.Compare(string(w1.Msn), string(w2.Msn)) == 0
}
