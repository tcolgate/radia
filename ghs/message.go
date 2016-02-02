package ghs

type Sender interface {
	Send(Message)
}

type Reciever interface {
	Recieve() Message
}

type Closer interface {
	Close()
}

type SenderReciever interface {
	Sender
	Reciever
	Closer
}

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

type chanPair struct {
	send chan<- Message
	recv <-chan Message
}

func (p chanPair) Send(m Message) {
	p.send <- m
}

func (p chanPair) Recieve() Message {
	return <-p.recv
}

func (p chanPair) Close() {
	close(p.send)
}

type SenderRecieverMaker func() (SenderReciever, SenderReciever)
