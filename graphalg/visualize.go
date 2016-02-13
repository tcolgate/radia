package graphalg

import (
	"fmt"
	"net/http"
	"text/template"

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

type Visualize struct {
	*http.ServeMux
	Nodes []*Node
	OnRun func()

	update chan struct{}
}

func (v Visualize) handleRoot(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	err := tmpl.Execute(w, struct{}{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func (v Visualize) updateSocket(ws *websocket.Conn) {
	for i := 0; i < 100; i++ {
		fmt.Fprintf(ws, "From here  %v", i)
		<-v.update
	}
}

func (v Visualize) handleRun(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	v.update <- struct{}{}
}

func MakeVisualize([]*Node) Visualize {
	v := Visualize{}
	v.update = make(chan struct{})
	v.ServeMux = http.NewServeMux()
	v.HandleFunc("/", v.handleRoot)
	v.Handle("/updates", websocket.Handler(v.updateSocket))
	v.HandleFunc("/run", v.handleRun)

	return v
}
