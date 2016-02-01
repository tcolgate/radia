package ghs

type FragmentID int

type Broadcast struct {
	FragID FragmentID
	Level  int
}

type Convergecast struct {
	MinWeight Weight
}

type MessageType int

const (
	MessageHello MessageType = iota
	MessageBroadcast
	MessageConvergeCast
)

type Message struct {
	Type MessageType
}
