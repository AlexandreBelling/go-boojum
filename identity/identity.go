package identity

import (
	"strings"
	"math/big"
	"io/ioutil"
	"crypto/ecdsa"
	"github.com/libp2p/go-libp2p-core/peer"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"

)

// PrivKey is a general purpose identity manager
//
// It is coercable in ethereum identities
// It is coercable in protocol members identities
// It is also coercable in ethereum private key
type PrivKey ecdsa.PrivateKey

// PubKey is a project local alias for ecdsa.PublicKey
// It is used to have coercion
type PubKey ecdsa.PublicKey

// ID is a generique alias for peer.ID compatible 
// It is usefull to avoid a dependency hell
type ID peer.ID

// Generate returns a generated private key
func Generate()	*PrivKey {
	res, _ := ethcrypto.GenerateKey()
	return (*PrivKey)(res)
}

// ReadPrivKey reads a private key from a file
func ReadPrivKey(pkPath string) (*PrivKey, error) {
	pkRaw, err := ioutil.ReadFile(pkPath)
	if err != nil {
		return nil, err
	}

	pkStr := string(pkRaw)
	pkStr = strings.TrimRight(pkStr, "\n")
	pkStr = strings.TrimRight(pkStr, "\r")

	pk, err := ethcrypto.HexToECDSA(pkStr)
	if err != nil {
		return nil, err
	}

	return (*PrivKey)(pk), nil
}

// PeerID returns the libp2p peer corresponding to the PrivKey
func (pk *PrivKey) PeerID() peer.ID {
	pid, _ := peer.IDFromPrivateKey(pk.Libp2p())
	return pid
}

// ID returns the libp2p peer corresponding to the PrivKey
func (pk *PrivKey) ID() ID {
	return ID(pk.PeerID())
}

// Libp2p returns the private key casted to be a p2pcrypto.PrivKey
func (pk *PrivKey) Libp2p() p2pcrypto.PrivKey {
	return (*p2pcrypto.Secp256k1PrivateKey)(pk)
}

// Eth returns the private key casted to be of geth type
func (pk *PrivKey) Eth() *ecdsa.PrivateKey {
	return (*ecdsa.PrivateKey)(pk)
}

// UnmarshalPubKey attemps to unmarshal a public key
func UnmarshalPubKey(b []byte) (*PubKey, error) {
	pub, err := ethcrypto.DecompressPubkey(b)
	if err != nil {
		return nil, err
	}
	return (*PubKey)(pub), nil
}

// Libp2p returns the public key casted to be a p2pcrypto.PrivKey
func (pub *PubKey) Libp2p() p2pcrypto.PubKey {
	return (*p2pcrypto.Secp256k1PublicKey)(pub)
}

// Eth returns the public key casted as an ethereum public key
func (pub *PubKey) Eth() *ecdsa.PublicKey {
	return (*ecdsa.PublicKey)(pub)
}

// PeerID returns the libp2p peer corresponding to the PubKey
func (pub *PubKey) PeerID() peer.ID {
	pid, _ := peer.IDFromPublicKey(pub.Libp2p())
	return pid
}

// ID returns the libp2p peer corresponding to the PubKey
func (pub *PubKey) ID() ID {
	return ID(pub.ID())
}

// Big returns a big.Int representation of the ID
func (id *ID) Big() *big.Int {
	x := big.NewInt(0)
	b := []byte(*id)
	x.SetBytes(b[:32])
	return x
}

func (id *ID) String() string {
	return string(*id)
}

// StringToID return an id with a string as input
func StringToID(s string) ID {
	return ID(s)
}