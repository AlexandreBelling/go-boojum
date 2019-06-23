package p2p

import (
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
)

// Identity ...
func Identity(sk crypto.PrivKey) (libp2p.Option, error) {
	return libp2p.Identity(sk), nil
}

// ListenAddress ...
func ListenAddress(addr string) (libp2p.Option, error) {
	return libp2p.ListenAddrStrings(addr), nil
}
