package graphalg

type Algorithm interface {
	Run()

	Edges() Edges
	Dispatcher
	Queuer
	Doner
}

type Queuer interface {
	Queue(Message)
	ClearQueue()
}

type Dispatcher interface {
	Dispatch(Message)
}

type Doner interface {
	Done() bool
	SetDone()
	OnDone() func()
	SetOnDone(func())
}
