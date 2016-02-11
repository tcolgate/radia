package graphalg

func Run(a Algorithm) {
	msgQueue := []Message{}
	ms := make(chan Message)
	defer close(ms)

	a.Edges().SortByMinEdge()
	defer func() {
		if a.OnDone() != nil {
			a.OnDone()
		}
	}()

	for _, e := range a.Edges() {
		go func(e *Edge) {
			for {
				ms <- e.Recieve()
			}
		}(e)
	}

	for nm := range ms {
		delayed := msgQueue
		msgQueue = []Message{}
		a.Dispatch(nm)
		for _, om := range delayed {
			a.Dispatch(om)
		}
		if a.Done() {
			return
		}
	}
}
