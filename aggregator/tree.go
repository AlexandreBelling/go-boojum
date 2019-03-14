package aggregator

import "C"

//Tree is a container
type Tree struct {
	boo *Boojum
	data *byte
}

func newTree(boo *Boojum) (*Tree) {
	var data *byte
	data = nil
	return &Tree{
		boo: boo,
		data : data,
	}
}

func (tree *Tree) Rm() {
	memFree(tree.data)
}

func (tree *Tree) Verify() (bool) {
	return verify(tree.data)
}
