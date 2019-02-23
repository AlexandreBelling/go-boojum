package aggregator

import (
	"testing"
)

func TestBoojum(t *testing.T) {

	boo := New().Initialize().WithDir("./setup").RunGenerators()

	left := boo.MakeExample()
	right := boo.MakeExample()

	output := boo.AggregateTrees(*left, *right)

	valid := boo.Verify(output)
	if !valid {
		t.FailNow()
	}

}