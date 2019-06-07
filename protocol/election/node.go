package election

import (
	"sync"

	"github.com/AlexandreBelling/go-boojum/protocol"
)

// Node contains data for each node of the aggregation tree 
type Node struct {
	Tree					*protocol.Tree
	Label					int
	AggregateProof			[]byte

	mut						sync.RWMutex
	readiness				int
	arity					int
	hookOnReadinessUpdate 	HookOnReadinessUpdate
}

// Job contains the 
type Job struct{
	InputProofs		[][]byte
	Label			int
}

// SetAggregateProof set the aggregation field and update the readiness of its parent
func (n *Node) SetAggregateProof(aggregateProof []byte) {
	n.AggregateProof = aggregateProof
	n.Tree.Parent.Node.(*Node).IncrementReadiness()
}

// IncrementReadiness update the readiness of the node
func (n *Node) IncrementReadiness() {
	defer n.hookOnReadinessUpdate(n)

	n.mut.Lock()
	defer n.mut.Unlock()
	n.readiness++
}

// IsReady returns true if the node can start being aggregated
func (n *Node) IsReady() bool {
	n.mut.RLock()
	defer n.mut.RUnlock()
	return n.readiness == n.arity
}

// Job provide a list of proofs to be aggregated
func (n *Node) Job() *Job {

	inputProofs := make([][]byte, n.arity)
	for index, children := range n.Tree.Children {
		inputProofs[index] = children.Node.(*Node).AggregateProof
	}

	return &Job{
		InputProofs: 	inputProofs,
		Label:			n.Label,
	}
}

// HookOnReadinessUpdate is function that is trigerred by the tree when we update its readiness
// We choose to apply a hook injection pattern because we don't implement control logic here
type HookOnReadinessUpdate func(n *Node)

// InitializeNodes performs all the election specific initialization
func InitializeNodes(t *protocol.Tree, f HookOnReadinessUpdate) {
	initializeNodes := makeNodeInitializer()
	applyOnReadinessUpdateHook := makeHookOnReadinessUpdateApplier(f)

	t.Walk(initializeNodes)
	t.Walk(applyOnReadinessUpdateHook)
}

// makeNodeInitializer returns a protocol.TreeFunc that initialize all the nodes of a tree
func makeNodeInitializer() protocol.TreeFunc {
	counter := protocol.MakeCounter() // Will be used to label the nodes with unique numbers
	return func(t *protocol.Tree) {

		arity := 0
		if t.Children != nil {
			arity = len(t.Children)
		}

		t.Node = &Node{
			Tree: 			t,
			Label:			counter(), // Gives a unique label to each node
			arity:			arity,
			readiness: 		0,
		}
	}
}

// makeHookOnReadinessUpdateApplier returns a TreeFunc setting the HookOnReadinessUpdate
func makeHookOnReadinessUpdateApplier(f HookOnReadinessUpdate) protocol.TreeFunc {
	return func(t *protocol.Tree) {
		t.Node.(*Node).hookOnReadinessUpdate = f
	}
}

// makeNodeMapperByLabel returns a map of node indexed by their label
func makeNodeMapperByLabel() (protocol.TreeFunc, func() map[int]*Node) {
	nodeMap := make(map[int]*Node)

	mapNode := func(t *protocol.Tree) { nodeMap[t.Node.(*Node).Label] = t.Node.(*Node) }
	getNodeMap := func() map[int]*Node { return nodeMap }

	return mapNode, getNodeMap
}
