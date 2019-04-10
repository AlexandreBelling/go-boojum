package layered

// Tree is a recursive helper that helps scheduling the aggregation
type Tree struct{
	children []*Tree
	payloadChan chan []byte
	payload []byte
	height int
}

// NewTree assigns a tree with depth 
func NewTree(height, arity int) (t* Tree) {
	if height > 0 {

		children := make([]*Tree, arity)
		for i:=0; i<arity; i++ {
			children[i] = NewTree(height-1, arity)
		}

		t = &Tree{
			children: children,
			payloadChan: make(chan []byte, 1),
			height: height,
		}

		return t
	}
	return nil
}

// GetLeaves returns the list of all leaves
func (t *Tree) GetLeaves() (leaves []Tree) {

	leaves = make ([]Tree, 0)
	t.getLeavesRecursive(&leaves)
	return leaves
}

func (t *Tree) getLeavesRecursive(leaves *[]Tree) {

	// This will work even if the tree is an incomplete tree
	for _, child := range t.children {
		if child != nil {
			child.getLeavesRecursive(leaves)
			return
		}
	}

	*leaves = append(*leaves, *t)
	return
}

