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
