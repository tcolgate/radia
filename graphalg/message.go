package graphalg

type Message struct {
	Edge  int
	Bytes []byte
}

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

type SenderRecieverMaker func() (SenderReciever, SenderReciever)
