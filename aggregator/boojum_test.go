package aggregator

import (
	"testing"
)

func TestBoojum(t *testing.T) {

	boo := NewBoojum().WithDir("./setup")

	left := boo.MakeExample()
	right := boo.MakeExample()

	boo.RunGenerators()
	output := boo.AggregateTrees(left, right)

	valid := boo.Verify(output)
	if !valid {
		t.FailNow()
	}

}
