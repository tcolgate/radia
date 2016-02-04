package ghs

type NodeID int

//go:generate stringer -type=NodeState
type NodeState int

const (
	NodeStateSleeping NodeState = iota
	NodeStateFind
	NodeStateFound
)

type Node struct {
	ID       NodeID
	State    NodeState
	Edges    Edges
	Level    int
	Fragment FragmentID

	bestEdge  int
	bestWt    int
	testEdge  int
	inBranch  int
	findCount int
}

func Join(n1 *Node, n2 *Node, f SenderRecieverMaker) {
	e1, e2 := NewEdge(f)
	n1.Edges = append(n1.Edges, e1)
	n2.Edges = append(n2.Edges, e2)
}

func (n1 *Node) Join(n2 *Node, f SenderRecieverMaker) {
	e1, e2 := NewEdge(f)
	n1.Edges = append(n1.Edges, e1)
	n2.Edges = append(n2.Edges, e2)
}

// WakeUp node procedure
func (n *Node) WakeUp() {
	e := n.Edges.MinEdge()
	n.State = NodeStateFound
	n.findCount = 0
	n.Level = 0
	e.SendConnect(n.Level)
}

// ChangeRoot node procedure
func (n *Node) ChangeRoot() {
}

func (n *Node) Connect(e *Edge, level int) {
}

func (n *Node) Initiate(e *Edge, level int, fragment FragmentID, state NodeState) {
}

func (n *Node) Test(e *Edge, level int, fragment FragmentID) {
}

func (n *Node) Accept(e *Edge) {
}

func (n *Node) Reject(e *Edge) {
}

func (n *Node) Report(e *Edge, best Weight) {
}
