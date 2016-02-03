package ghs

import (
	"log"
	"testing"
)

func TestEdgeTestMessage1(t *testing.T) {
	n1 := Node{}
	n2 := Node{}
	Join(&n1, &n2, MakeChanPair)

	if len(n1.Edges) < 1 {
		t.Fatalf("expected %v edges, got %v", 1, len(n1.Edges))
	}
	if len(n2.Edges) < 1 {
		t.Fatalf("expected %v edges, got %v", 1, len(n2.Edges))
	}

	go func() { n1.Edges[0].Send(sentm) }()
	gotm := n2.Edges[0].Recieve()

	log.Printf("%+v, %+v", sentm, gotm)

	if sentm != gotm {
		t.Fatalf("expected %+v, got %+v", 1, len(n2.Edges))
	}
}

func TestEdgeTestMessage2(t *testing.T) {
	n1 := Node{}
	n2 := Node{}
	n1.Join(&n2, MakeChanPair)
	//m := Message{}
}

func TestEdgesNextEmpty(t *testing.T) {

}
