package boojum

//Tree is a container
type Tree struct {
	data *[]byte
}

func newTree() (*Tree) {
	data := make([]byte, 1)
	return &Tree{
		data : &data,
	}
}

