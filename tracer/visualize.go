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

package tracer

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

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
	T        int64
	NodeID   string
	EdgeName string `json:",omitempty"`
	Log      string `json:",omitempty"`
	State    string `json:",omitempty"`
	Message  string `json:",omitempty"`
}

type httpDisplay struct {
	OnRun  func()
	update chan struct{}
	msgs   chan jsonMsg
}

func (h *httpDisplay) Log(t int64, id, s string) {
	h.msgs <- jsonMsg{T: t, NodeID: id, Log: s}
}

func (h *httpDisplay) NodeUpdate(t int64, n, str string) {
	h.msgs <- jsonMsg{T: t, NodeID: n, State: str}
}

func (h *httpDisplay) EdgeUpdate(t int64, n, en, s string) {
	h.msgs <- jsonMsg{T: t, NodeID: n, EdgeName: en, State: s}
}

func (h *httpDisplay) EdgeMessage(t int64, n, en, m string) {
	h.msgs <- jsonMsg{T: t, NodeID: n, EdgeName: en, Message: m}
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
	type jn struct {
		Id string `json:"id"`
	}
	type jl struct {
		Source int     `json:"source"`
		Target int     `json:"target"`
		Cost   float64 `json:"cost"`
	}
	type d struct {
		Nodes []jn `json:"nodes"`
		Links []jl `json:"links"`
	}

	for {
		select {
		case m := <-v.msgs:
			log.Println(m)
			b, _ := json.Marshal(m)
			log.Println(b)
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

func NewHTTPDisplay(mux *http.ServeMux, onRun func()) *Tracer {
	v := httpDisplay{}
	v.OnRun = onRun
	v.update = make(chan struct{})
	v.msgs = make(chan jsonMsg)

	mux.HandleFunc("/", v.handleRoot)
	mux.Handle("/updates", websocket.Handler(v.updateSocket))
	mux.HandleFunc("/run", v.handleRun)

	return &Tracer{&v}
}
