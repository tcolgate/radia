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

	bestEdge  *Edge
	bestWt    Weight
	testEdge  *Edge
	inBranch  *Edge
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

// procWakeUp - wakeup node procedure
func (n *Node) procWakeUp() {
	e := n.Edges.MinEdge()
	n.State = NodeStateFound
	n.findCount = 0
	n.Level = 0
	e.SendConnect(n.Level)
}

// ProcChangeRoot - change-rootnode procedure
func (n *Node) procChangeRoot() {
}

// ProcTest - test node procedure
func (n *Node) procTest() {
}

// ProcReport - report node procedure
func (n *Node) procReport() {
}

func (n *Node) Connect(e *Edge, level int) {
	if n.State == NodeStateSleeping {
		n.procWakeUp()
	}
	if level < n.Level {
		e.State = EdgeStateBranch
		e.SendInitiate(n.Level, n.Fragment, n.State)
		if n.State == NodeStateFind {
			n.findCount++
		}
	} else if e.State == EdgeStateBasic {
		// Need to queue this
	} else {
		e.SendInitiate(n.Level+1, e.Weight, NodeStateFind)
	}
}

func (n *Node) Initiate(e *Edge, level int, fragment FragmentID, state NodeState) {
	n.Level = level
	n.Fragment = fragment
	n.State = state
	n.inBranch = e
	n.bestEdge = nil
	n.bestWt = WeightInf
	for _, se := range n.Edges {
		if se != e && e.State == EdgeStateBranch {
			se.SendInitiate(n.Level, n.Fragment, n.State)
			if n.State == NodeStateFind {
				n.findCount++
			}
		}
		if n.State == NodeStateFind {
			n.procTest()
		}
	}
}

func (n *Node) Test(e *Edge, level int, fragment FragmentID) {
}

func (n *Node) Accept(e *Edge) {
}

func (n *Node) Reject(e *Edge) {
}

func (n *Node) Report(e *Edge, best Weight) {
}
