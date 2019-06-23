package protocol

import (
	"github.com/AlexandreBelling/go-boojum/identity"
	"github.com/ethereum/go-ethereum/crypto"
)

// IDFromBatch gets an ID from the hash of a protocol
func IDFromBatch(batch [][]byte) identity.ID {
	h := crypto.Keccak256Hash(batch...).String()
	return identity.ID(h)
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
