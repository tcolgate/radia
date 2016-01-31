package ghs

type FragmentID int

type BroadcastMsg struct {
	FragID FragmentID
	Level  int
}

type ConvergecastMsg struct {
	MinWeight Weight
}
