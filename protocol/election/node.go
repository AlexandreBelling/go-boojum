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

	mut						sync.Mutex
	readiness				int
	arity					int
	hookOnReadinessUpdate 	HookOnReadinessUpdate
}

// HookOnReadinessUpdate is function that is trigerred by the tree when we update its readiness
// We choose to apply a hook injection pattern because we don't implement control logic here
type HookOnReadinessUpdate func(n *Node)

// MakeNodeInitializer returns a protocol.TreeFunc that initialize all the nodes of a tree
func MakeNodeInitializer() protocol.TreeFunc {
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

// MakeNodeMapperByLabel returns a map of node indexed by their label
func MakeNodeMapperByLabel() (protocol.TreeFunc, func() map[int]*Node) {
	nodeMap := make(map[int]*Node)

	mapNode := func(t *protocol.Tree) { nodeMap[t.Node.(*Node).Label] = t.Node.(*Node) }
	getNodeMap := func() map[int]*Node { return nodeMap }

	return mapNode, getNodeMap
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

// Set 
