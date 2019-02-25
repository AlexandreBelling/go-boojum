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
