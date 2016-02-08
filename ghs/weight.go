package ghs

import "math"

type Weight struct {
	float64
	Lsn NodeID // Least Signigicant NodeID
	Msn NodeID // Most Signigicant NodeID
}

var WeightInf = Weight{float64: math.Inf(1)}

// Less - compare two edge weights.
func (w1 Weight) Less(w2 Weight) bool {
	switch {
	case w1.float64 < w2.float64,
		w1.float64 == w2.float64 && int(w1.Lsn) < int(w2.Lsn),
		w1.float64 == w2.float64 && int(w1.Lsn) == int(w2.Lsn) && int(w1.Msn) < int(w2.Msn):
		return true
	default:
		return false
	}
}

// Less - compare two edge weights.
func (w1 Weight) Greater(w2 Weight) bool {
	switch {
	case w1.float64 > w2.float64,
		w1.float64 == w2.float64 && int(w1.Lsn) > int(w2.Lsn),
		w1.float64 == w2.float64 && int(w1.Lsn) == int(w2.Lsn) && int(w1.Msn) > int(w2.Msn):
		return true
	default:
		return false
	}
}

// Less - compare two edge weights.
func (w1 Weight) Equal(w2 Weight) bool {
	return w1.float64 == w2.float64 && int(w1.Lsn) == int(w2.Lsn) && int(w1.Msn) == int(w2.Msn)
}

// FragmentID converts a Weight to a FragmentID. The details of the best
//  edge in a fragment are effectively act as a fragment id.
// In 2) Response to receipt of Connect(... we have
// ...
//    else send Initiate(LN + 1, w(j), Find) on edge j
// ...
// Which clearly uses the edge weight in a Initiate (L, F, S)  message
func (w Weight) FragmentID() FragmentID {
	return FragmentID(w)
}
