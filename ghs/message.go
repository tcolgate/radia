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

//go:generate stringer -type=MessageType
type MessageType int

const (
	MessageConnect MessageType = iota
	MessageInitiate
	MessageTest
	MessageAccept
	MessageReject
	MessageReport
	MessageChangeRoot
)

type Message struct {
	Edge *Edge
	Type MessageType
	Weight
	FragmentID
	Level int
}

type SenderRecieverMaker func() (SenderReciever, SenderReciever)

type chanPair struct {
	send chan<- Message
	recv <-chan Message
}

func MakeChanPair() (SenderReciever, SenderReciever) {
	c1, c2 := make(chan Message), make(chan Message)
	return chanPair{c1, c2}, chanPair{c2, c1}
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
