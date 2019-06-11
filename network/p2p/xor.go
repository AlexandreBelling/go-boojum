package p2p

import (
	"math/big"
	"sort"

	"github.com/libp2p/go-libp2p-core/peer"
)

// ByXORDistance is util struct for sorting peers by their XOR distance between each others
type byXORDistance struct {
	slice     []peer.ID
	reference peer.ID
}

func (x byXORDistance) Len() int      { return len(x.slice) }
func (x byXORDistance) Swap(i, j int) { x.slice[i], x.slice[j] = x.slice[j], x.slice[i] }
func (x byXORDistance) Less(i, j int) bool {
	left := XORDistance(x.reference, x.slice[i])
	right := XORDistance(x.reference, x.slice[j])
	return left.Cmp(&right) == -1
}
func (x byXORDistance) Sort() { sort.Sort(x) }

// XORDistance returns the XOR distance between two peers
func XORDistance(a, b peer.ID) big.Int {

	aBin := []byte(a)
	bBin := []byte(b)
	cBin := make([]byte, len(aBin))

	for i := range aBin {
		cBin[i] = aBin[i] ^ bBin[i]
	}

	var distance big.Int
	distance.SetBytes(cBin)
	return distance
}

// Sort the
