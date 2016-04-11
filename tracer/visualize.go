// Copyright (c) 2016 Tristan Colgate-McFarlane
//
// This file is part of radia.
//
// radia is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// radia is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with radia.  If not, see <http://www.gnu.org/licenses/>.

package tracer

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/tcolgate/radia/graph"

	"golang.org/x/net/websocket"
)

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go assets/
func init() {
	data, err := Asset("assets/index.html.tmpl")
	if err != nil {
		panic(err.Error())
	}
	tmpl = template.Must(template.New("index.html.tmpl").Parse(string(data)))
}

var tmpl *template.Template

type jsonMsg struct {
	T        time.Time
	Type     string
	NodeID   string
	EdgeName string `json:",omitempty"`
	Log      string `json:",omitempty"`
	State    string `json:",omitempty"`
	Message  string `json:",omitempty"`
	Dir      string `json:",omitempty"`
}

type httpDisplay struct {
	OnRun  func()
	update chan struct{}
	msgs   chan jsonMsg
}

func (h *httpDisplay) Log(t time.Time, gid graph.GraphID, aid graph.AlgorithmID, id, s string) {
	h.msgs <- jsonMsg{T: t, Type: "log", NodeID: id, Log: s}
}

func (h *httpDisplay) NodeUpdate(t time.Time, gid graph.GraphID, aid graph.AlgorithmID, n, str string) {
	h.msgs <- jsonMsg{T: t, Type: "nodeUpdate", NodeID: n, State: str}
}

func (h *httpDisplay) EdgeUpdate(t time.Time, gid graph.GraphID, aid graph.AlgorithmID, n, en, s string) {
	h.msgs <- jsonMsg{T: t, Type: "edgeUpdate", NodeID: n, EdgeName: en, State: s}
}

func (h *httpDisplay) EdgeMessage(t time.Time, gid graph.GraphID, aid graph.AlgorithmID, n, en string, dir MessageDir, m string) {
	h.msgs <- jsonMsg{T: t, Type: "edgeMessage", NodeID: n, EdgeName: en, Dir: dir.String(), Message: m}
}

func (v httpDisplay) handleRoot(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	err := tmpl.Execute(w, struct{}{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func (v httpDisplay) updateSocket(ws *websocket.Conn) {
	for {
		select {
		case m := <-v.msgs:
			b, _ := json.Marshal(m)
			fmt.Fprintf(ws, string(b))
		}
	}

	/*

		jns := []jn{}
		jls := []jl{}

			nix := map[string]int{}

			//	for i, n := range v.Nodes {
			//		jns = append(jns, jn{string(n.ID)})
			//		nix[string(n.ID)] = i
			//	}

			data := d{jns, jls}
			b, _ := json.Marshal(data)

			fmt.Fprintf(ws, string(b))

			for {
				jns := []jn{}
				jls := []jl{}
				type eIndex struct{ l, r string }
					eix := map[eIndex]bool{}

						for _, n := range v.Nodes {
							for _, e := range n.Edges() {
								ei := eIndex{e.Weight.LsnID, e.Weight.MsnID}
								if _, ok := eix[ei]; !ok && !e.Disabled {
									jls = append(jls, jl{
										Source: nix[string(e.Weight.LsnID)],
										Target: nix[string(e.Weight.MsnID)],
										Cost:   e.Weight.Cost,
									})
									eix[ei] = true
								}
							}
						}

				data := d{jns, jls}
				b, _ := json.Marshal(data)

				fmt.Fprintf(ws, string(b))
				<-v.update
			}
	*/
}

func (v httpDisplay) handleRun(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	v.OnRun()
	v.update <- struct{}{}
}

func NewHTTPDisplay(mux *http.ServeMux, onRun func()) traceDisplay {
	v := httpDisplay{}
	v.OnRun = onRun
	v.update = make(chan struct{})
	v.msgs = make(chan jsonMsg)

	mux.HandleFunc("/", v.handleRoot)
	mux.Handle("/updates", websocket.Handler(v.updateSocket))
	mux.HandleFunc("/run", v.handleRun)

	return &v
}
