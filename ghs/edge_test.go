package ghs

import "testing"

func TestEdgeTestMessage1(t *testing.T) {
	n1 := Node{}
	n2 := Node{}
	Join(&n1, &n2, MakeChanPair)

	if len(n1.Edges) != 1 {
		t.Fatalf("expected %v edges, got %v", 1, len(n1.Edges))
	}
	if len(n2.Edges) != 1 {
		t.Fatalf("expected %v edges, got %v", 1, len(n2.Edges))
	}

	sentm := Message{Type: MessageConnect}
	go func() { n1.Edges[0].Send(sentm) }()
	gotm := n2.Edges[0].Recieve()

	gotm.Edge = nil
	if sentm != gotm {
		t.Fatalf("expected %+v, got %+v", sentm, gotm)
	}
}

func TestEdgeTestMessage2(t *testing.T) {
	n1 := Node{}
	n2 := Node{}
	n3 := Node{}
	Join(&n1, &n2, MakeChanPair)
	Join(&n3, &n2, MakeChanPair)

	if len(n1.Edges) != 1 {
		t.Fatalf("expected %v edges, got %v", 1, len(n1.Edges))
	}
	if len(n2.Edges) != 2 {
		t.Fatalf("expected %v edges, got %v", 1, len(n2.Edges))
	}

	if len(n3.Edges) != 1 {
		t.Fatalf("expected %v edges, got %v", 1, len(n2.Edges))
	}

	sentm := Message{Type: MessageConnect}
	c := make(chan Message)
	go func() { c <- n1.Edges[0].Recieve() }()
	go func() { n3.Edges[0].Send(sentm) }()
	gotm := n2.Edges[1].Recieve()
	select {
	case <-c:
		t.Fatalf("expected recieve to block")
	default:
	}

	gotm.Edge = nil
	if sentm != gotm {
		t.Fatalf("expected %+v, got %+v", sentm, gotm)
	}
}

func TestEdgesNextEmpty(t *testing.T) {

}
