package ghs

type FragmentID int

type Broadcast struct {
	FragID FragmentID
	Level  int
}

type Convergecast struct {
	MinWeight Weight
}

type Message struct {
	*Broadcast
	*Convergecast
}
