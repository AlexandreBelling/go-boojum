package protocol

import (
	"math/big"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"
)

// ID is a general purpose identifier for identities
type ID common.Hash

// Big returns a big.Int representation of the ID
func (id *ID) Big() *big.Int {
	return (*common.Hash)(id).Big()
}

// IDFromBatch gets an ID from the hash of a protocol
func IDFromBatch(batch [][]byte) ID {
	h := crypto.Keccak256Hash(batch...)
	return ID(h)
}

// MakeCounter return a simple basic counter
func MakeCounter() Counter {
	value := 0
	return func() int {
		value++
		return value
	}
}

// Counter is a function that return an incremented value at each call by one
type Counter func() int