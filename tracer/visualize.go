package tracer

//go:generate protoc -I $GOPATH/src:. --js_out=assets/ internal/proto/tracer.proto
//go:generate go-bindata -pkg $GOPACKAGE -o assets.go assets/
/*
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
	//Nodes []*graphalg.Node
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
}

func (v Visualize) handleRun(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	v.OnRun()
	v.update <- struct{}{}
}

func MakeVisualize(ns []*Node, onRun func()) Visualize {
	v := Visualize{}
	v.Nodes = ns
	v.OnRun = onRun
	v.update = make(chan struct{})
	v.ServeMux = http.NewServeMux()
	v.HandleFunc("/", v.handleRoot)
	v.Handle("/updates", websocket.Handler(v.updateSocket))
	v.HandleFunc("/run", v.handleRun)

	return v
}
*/
