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

package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/tcolgate/vonq/graphalg"
)

func setupGHS() graphalg.Visualize {
	wg := sync.WaitGroup{}

	// We'll only ever get halt messages from the core edge, so only
	// two nodes halt
	wg.Add(2)
	n1 := graphalg.Node{
		ID:     graphalg.NodeID("n1"),
		Logger: log.New(os.Stdout, "node(n1)", 0),
	}
	n2 := graphalg.Node{
		ID:     graphalg.NodeID("n2"),
		Logger: log.New(os.Stdout, "node(n2)", 0),
	}
	n3 := graphalg.Node{
		ID:     graphalg.NodeID("n3"),
		Logger: log.New(os.Stdout, "node(n3)", 0),
	}
	n4 := graphalg.Node{
		ID:     graphalg.NodeID("n4"),
		Logger: log.New(os.Stdout, "node(n4)", 0),
	}

	n1.Join(&n2, 1.0, graphalg.MakeChanPair)
	n2.Join(&n3, 1.0, graphalg.MakeChanPair)
	n3.Join(&n1, 1.0, graphalg.MakeChanPair)
	n1.Join(&n4, 1.0, graphalg.MakeChanPair)

	/*
		ghs1 := ghs.State{Node: &n1}
		ghs2 := ghs.State{Node: &n2}
		ghs3 := ghs.State{Node: &n3}
		ghs4 := ghs.State{Node: &n4}

			go n1.Run(&ghs1, wg.Done)
			go n2.Run(&ghs2, wg.Done)
			go n3.Run(&ghs3, wg.Done)
			go n4.Run(&ghs4, wg.Done)

			ghs1.WakeUp()
			wg.Wait()
	*/
	nodes := []*graphalg.Node{
		&n1,
		&n2,
		&n3,
		&n4,
	}

	return graphalg.MakeVisualize(nodes)
}

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func main() {
	err := http.ListenAndServe(":12345", setupGHS())
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
