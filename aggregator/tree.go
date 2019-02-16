package boojum

type Tree struct {
	data *[]byte
}

func newTree() (*Tree) {
	return &Tree{
		data : make(*[]byte)
	}
}

