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
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/tcolgate/vonq/ghs"
	"github.com/tcolgate/vonq/graphalg"
	"github.com/tcolgate/vonq/tracer"
)

type state struct {
	ghs []*ghs.State
	wg  *sync.WaitGroup
}

func (s *state) OnRun() {
	s.ghs[0].WakeUp()
	s.wg.Wait()

	log.Println("finished")
}

func setupGHS(mux *http.ServeMux) {
	s := state{}

	t := tracer.NewHTTPDisplay(mux, s.OnRun)

	s.wg = &sync.WaitGroup{}

	count := 25
	s.wg.Add(count)
	nodes := make([]*graphalg.Node, 0)
	s.ghs = make([]*ghs.State, 0)

	for i := 0; i < count; i++ {
		n := &graphalg.Node{
			ID:     graphalg.NodeID(fmt.Sprintf("n%v", i+1)),
			Tracer: t,
		}
		nodes = append(nodes, n)
		s.ghs = append(s.ghs, &ghs.State{Node: n})
	}

	for i := 0; i < count; i++ {
		for j := i + 1; j < count; j++ {
			log.Println("join", i, j)
			w := rand.NormFloat64()*1.0 + 4.0
			nodes[i].Join(nodes[j], w, graphalg.MakeChanPair)
		}
	}

	for _, g := range s.ghs {
		go func(g *ghs.State) {
			n := g.Node
			n.NodeUpdate("from node")
			for _, e := range n.Edges() {
				e.EdgeUpdate()
			}

			graphalg.Run(g, s.wg.Done)
		}(g)
	}

	log.Println("running")
}

func main() {
	mux := http.NewServeMux()
	setupGHS(mux)
	err := http.ListenAndServe(":12345", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
