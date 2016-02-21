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
	Log string
}

type httpDisplay struct {
	OnRun  func()
	update chan struct{}
	msgs   chan jsonMsg
}

func (h *httpDisplay) Log(s string) {
}

func (h *httpDisplay) NodeUpdate() {
}

func (h *httpDisplay) EdgeUpdate() {
}

func (h *httpDisplay) EdgeMessage(str string) {
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
