package ghs

import (
	"log"
	"os"
	"sync"
	"testing"

	"golang.org/x/net/context"
)

func TestGHSTest1(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(3)

	n1 := Node{
		ID:       NodeID("n1"),
		OnDone:   wg.Done,
		Fragment: FragmentID{Msn: NodeID("n1")},
		Logger:   log.New(os.Stdout, "node(n1) ", 0),
	}
	n2 := Node{
		ID:       NodeID("n2"),
		OnDone:   wg.Done,
		Fragment: FragmentID{Msn: NodeID("n2")},
		Logger:   log.New(os.Stdout, "node(n2) ", 0),
	}

	Join(&n1, &n2, 1.0, MakeChanPair)

	go n1.Run(ctx)
	go n2.Run(ctx)

	n1.WakeUp(ctx)

	wg.Wait()
	cancel()
}

/*
func TestGHSTest2(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(3)

	n1 := Node{
		ID:       NodeID("n1"),
		OnDone:   wg.Done,
		Fragment: FragmentID{Msn: NodeID("n1")},
		Logger:   log.New(os.Stdout, "node(n1) ", 0),
	}
	n2 := Node{
		ID:       NodeID("n2"),
		OnDone:   wg.Done,
		Fragment: FragmentID{Msn: NodeID("n2")},
		Logger:   log.New(os.Stdout, "node(n2) ", 0),
	}

	n3 := Node{
		ID:       NodeID("n3"),
		OnDone:   wg.Done,
		Fragment: FragmentID{Msn: NodeID("n3")},
		Logger:   log.New(os.Stdout, "node(n3) ", 0),
	}

	Join(&n1, &n2, 1.0, MakeChanPair)
	Join(&n3, &n2, 1.1, MakeChanPair)

	go n1.Run()
	go n2.Run()
	go n3.Run()

	n1.WakeUp()

	wg.Wait()
}

func TestGHSTest3(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(3)

	n1 := Node{
		ID:       NodeID("n1"),
		OnDone:   wg.Done,
		Fragment: FragmentID{float64: 10, Msn: NodeID("n1")},
		Logger:   log.New(os.Stdout, "node(n1) ", 0),
	}
	n2 := Node{
		ID:       NodeID("n2"),
		OnDone:   wg.Done,
		Fragment: FragmentID{float64: 12, Msn: NodeID("n2")},
		Logger:   log.New(os.Stdout, "node(n2) ", 0),
	}

	n3 := Node{
		ID:       NodeID("n3"),
		OnDone:   wg.Done,
		Fragment: FragmentID{float64: 11, Msn: NodeID("n3")},
		Logger:   log.New(os.Stdout, "node(n3) ", 0),
	}

	Join(&n1, &n2, 1.0, MakeChanPair)
	Join(&n3, &n2, 1.1, MakeChanPair)
	Join(&n3, &n1, 2.1, MakeChanPair)

	go n1.Run()
	go n2.Run()
	go n3.Run()

	n1.WakeUp()

	wg.Wait()
}
*/
