package aggregator

import "C"

//Tree is a container
type Tree struct {
	data *byte
}

func newTree() (*Tree) {
	var data *byte
	return &Tree{
		data : data,
	}
}

// Free applies the custom free function on the remaining elements
func (t *Tree) Free() {
	memFreeTree(t.data)
}

