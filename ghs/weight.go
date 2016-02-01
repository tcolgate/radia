package ghs

type Weight struct {
	float64
	Lsn NodeID // Least Signigicant NodeID
	Msn NodeID // Most Signigicant NodeID
}

type Weights []*Weight

func (w Weights) Len() int { return len(w) }

func (w Weights) Swap(i, j int) { w[i], w[j] = w[j], w[i] }

// Less - compare two edge weights. If we can make all the edges in the graph
// have unqiue weights, we should be guaranteed one unique topology
func (w Weights) Less(i, j int) bool {
	switch {
	case w[i].float64 < w[j].float64,
		w[i].float64 == w[j].float64 && int(w[i].Lsn) < int(w[j].Lsn),
		w[i].float64 == w[j].float64 && int(w[i].Lsn) == int(w[j].Lsn) && int(w[i].Msn) < int(w[j].Msn):
		return true
	default:
		return false
	}
}
