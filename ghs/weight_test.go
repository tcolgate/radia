package ghs

import "testing"

func TestWeigjt(t *testing.T) {
	w0 := WeightInf
	w1 := Weight{1.0, "n1", "n2"}
	w2 := Weight{1.1, "n2", "n3"}
	w3 := Weight{1.2, "n2", "n3"}

	if !w0.Greater(w1) {
		t.Fatal("! w0 > w1")
	}
	if !w1.Less(w2) {
		t.Fatal("! w1 < w2")
	}
	if !w2.Less(w3) {
		t.Fatal("! w2 < w3")
	}
	if w0.Less(w1) {
		t.Fatal("! w0 < w1")
	}
	if w1.Greater(w2) {
		t.Fatal("w1 > w2")
	}
	if w2.Greater(w3) {
		t.Fatal("w2 > w3")
	}
}
