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
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/tcolgate/vonq/ghs"
	"github.com/tcolgate/vonq/graphalg"
	"github.com/tcolgate/vonq/tracer"
)

type thing struct{}

type (*thing)OnRun(){
}

func setupGHS() *http.ServeMux {
	t := tracer.NewHTTPDisplay(nil)
	wg := sync.WaitGroup{}

	// We'll only ever get halt messages from the core edge, so only
	// two nodes halt
	wg.Add(2)
	n1 := graphalg.Node{
		ID:     graphalg.NodeID("n1"),
		Tracer: t,
	}
	n2 := graphalg.Node{
		ID:     graphalg.NodeID("n2"),
		Tracer: t,
	}
	n3 := graphalg.Node{
		ID:     graphalg.NodeID("n3"),
		Tracer: t,
	}
	n4 := graphalg.Node{
		ID:     graphalg.NodeID("n4"),
		Tracer: t,
	}
	n5 := graphalg.Node{
		ID:     graphalg.NodeID("n5"),
		Tracer: t,
	}
	n6 := graphalg.Node{
		ID:     graphalg.NodeID("n6"),
		Tracer: t,
	}

	n1.Join(&n2, 1.1, graphalg.MakeChanPair)
	n2.Join(&n4, 3.1, graphalg.MakeChanPair)
	n4.Join(&n6, 3.7, graphalg.MakeChanPair)
	n6.Join(&n5, 2.1, graphalg.MakeChanPair)
	n5.Join(&n3, 3.8, graphalg.MakeChanPair)
	n5.Join(&n1, 2.6, graphalg.MakeChanPair)
	n3.Join(&n1, 1.7, graphalg.MakeChanPair)

	ghs1 := ghs.State{Node: &n1}
	ghs2 := ghs.State{Node: &n2}
	ghs3 := ghs.State{Node: &n3}
	ghs4 := ghs.State{Node: &n4}
	ghs5 := ghs.State{Node: &n5}
	ghs6 := ghs.State{Node: &n6}

	onRun := func() {
		go n1.Run(&ghs1, wg.Done)
		go n2.Run(&ghs2, wg.Done)
		go n3.Run(&ghs3, wg.Done)
		go n4.Run(&ghs4, wg.Done)
		go n5.Run(&ghs5, wg.Done)
		go n6.Run(&ghs6, wg.Done)

		ghs1.WakeUp()
		wg.Wait()

		b, _ := json.Marshal(ghs6)

		log.Println(string(b))
		log.Println("finished")
	}

	nodes := []*graphalg.Node{
		&n1,
		&n2,
		&n3,
		&n4,
		&n5,
		&n6,
	}

	return t
}

func main() {
	err := http.ListenAndServe(":12345", setupGHS())
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
