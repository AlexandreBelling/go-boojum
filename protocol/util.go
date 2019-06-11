package protocol

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

// ID is a general purpose identifier for identities
type ID common.Hash

// Big returns a big.Int representation of the ID
func (id *ID) Big() *big.Int {
	return (*common.Hash)(id).Big()
}

func (id *ID) String() string {
	return (*common.Hash)(id).Hex()
}

// StringToID return an id with a string as input
func StringToID(s string) ID {
	return ID(common.HexToHash(s))
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
