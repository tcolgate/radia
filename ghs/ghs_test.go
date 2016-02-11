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

import (
	"log"
	"os"
	"sync"
	"testing"
)

func TestGHS1(t *testing.T) {
	wg := sync.WaitGroup{}
	// We'll only ever get halt messages from the core edge, so only
	// two nodes halt
	wg.Add(2)
	n1 := Node{
		ID:       NodeID("n1"),
		Logger:   log.New(os.Stdout, "node(n1)", 0),
		OnDone:   wg.Done,
		Fragment: FragmentID{float64: 1, Lsn: "blah", Msn: "n1"},
	}
	n2 := Node{
		ID:       NodeID("n2"),
		Logger:   log.New(os.Stdout, "node(n2)", 0),
		OnDone:   wg.Done,
		Fragment: FragmentID{float64: 2, Lsn: "blab", Msn: "n4"},
	}

	n1.Join(&n2, 1.0, MakeChanPair)

	go n1.Run()
	go n2.Run()

	n1.WakeUp()

	wg.Wait()
}

func TestGHS2(t *testing.T) {
	wg := sync.WaitGroup{}
	// We'll only ever get halt messages from the core edge, so only
	// two nodes halt
	wg.Add(2)
	n1 := Node{
		ID:       NodeID("n1"),
		Logger:   log.New(os.Stdout, "node(n1)", 0),
		OnDone:   wg.Done,
		Fragment: FragmentID{Msn: "n1"},
	}
	n2 := Node{
		ID:       NodeID("n2"),
		Logger:   log.New(os.Stdout, "node(n2)", 0),
		OnDone:   wg.Done,
		Fragment: FragmentID{Msn: "n2"},
	}
	n3 := Node{
		ID:       NodeID("n3"),
		Logger:   log.New(os.Stdout, "node(n3)", 0),
		OnDone:   wg.Done,
		Fragment: FragmentID{Msn: "n3"},
	}

	n1.Join(&n2, 1.0, MakeChanPair)
	n2.Join(&n3, 2.0, MakeChanPair)

	go n1.Run()
	go n2.Run()
	go n3.Run()

	n1.WakeUp()

	wg.Wait()
}

func TestGHS3(t *testing.T) {
	wg := sync.WaitGroup{}
	// We'll only ever get halt messages from the core edge, so only
	// two nodes halt
	wg.Add(2)
	n1 := Node{
		ID:       NodeID("n1"),
		Logger:   log.New(os.Stdout, "node(n1)", 0),
		OnDone:   wg.Done,
		Fragment: FragmentID{Msn: "n1"},
	}
	n2 := Node{
		ID:       NodeID("n2"),
		Logger:   log.New(os.Stdout, "node(n2)", 0),
		OnDone:   wg.Done,
		Fragment: FragmentID{Msn: "n2"},
	}
	n3 := Node{
		ID:       NodeID("n3"),
		Logger:   log.New(os.Stdout, "node(n3)", 0),
		OnDone:   wg.Done,
		Fragment: FragmentID{Msn: "n3"},
	}

	n1.Join(&n2, 1.0, MakeChanPair)
	n2.Join(&n3, 2.0, MakeChanPair)
	n3.Join(&n1, 3.0, MakeChanPair)

	go n1.Run()
	go n2.Run()
	go n3.Run()

	n1.WakeUp()

	wg.Wait()
}

func TestGHS4(t *testing.T) {
	wg := sync.WaitGroup{}
	// We'll only ever get halt messages from the core edge, so only
	// two nodes halt
	wg.Add(2)
	n1 := Node{
		ID:       NodeID("n1"),
		Logger:   log.New(os.Stdout, "node(n1)", 0),
		OnDone:   wg.Done,
		Fragment: FragmentID{Msn: "n1"},
	}
	n2 := Node{
		ID:       NodeID("n2"),
		Logger:   log.New(os.Stdout, "node(n2)", 0),
		OnDone:   wg.Done,
		Fragment: FragmentID{Msn: "n2"},
	}
	n3 := Node{
		ID:       NodeID("n3"),
		Logger:   log.New(os.Stdout, "node(n3)", 0),
		OnDone:   wg.Done,
		Fragment: FragmentID{Msn: "n3"},
	}
	n4 := Node{
		ID:       NodeID("n4"),
		Logger:   log.New(os.Stdout, "node(n4)", 0),
		OnDone:   wg.Done,
		Fragment: FragmentID{Msn: "n4"},
	}

	n1.Join(&n2, 1.0, MakeChanPair)
	n2.Join(&n3, 2.0, MakeChanPair)
	n3.Join(&n1, 3.0, MakeChanPair)
	n1.Join(&n4, 0.5, MakeChanPair)

	go n1.Run()
	go n2.Run()
	go n3.Run()
	go n4.Run()

	n1.WakeUp()

	wg.Wait()
}
