package ethereum

import (
	"sync"
	"context"
	"math/big"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/core/types"
)

type signerType string

const (
	// Homestead is an identifier to be used to specify 
	// the sender we want to use homestead signature
	Homestead 	signerType = "Homestead"
	// EIP155 is an identifier to be used to specify 
	// the sender we want to use EIP155
	EIP155		signerType = "EIP155"
)

// Sender can be provided a private key and use it to sign
//
// It is responsible for adding
type Sender struct {
	ctx 		context.Context
	cancel		context.CancelFunc

	priv 		ecdsa.PrivateKey
	address		common.Address

	client		*ethclient.Client
	chainID		*big.Int
	signer		types.Signer

	nonce		uint64
	mutex		sync.Mutex
}

// NewSender returns a Sender object
//
// It does not include any specific logic aside deriving the adress from the private key
func NewSender(
	ctx context.Context, 
	priv ecdsa.PrivateKey, 
	client *ethclient.Client,
	signerT signerType, 
	chainID *big.Int,
	) *Sender {

	ctx, cancel := context.WithCancel(ctx)
	pub := priv.Public().(ecdsa.PublicKey)

	sender := &Sender{
		ctx:		ctx,
		cancel: 	cancel,

		priv:		priv,
		address:	crypto.PubkeyToAddress(pub),

		client:		client,
		chainID:	chainID,
	}

	switch signerT {
	case EIP155:
		sender.signer = types.NewEIP155Signer(chainID)
	case Homestead:
		sender.signer = types.HomesteadSigner{}
	}
	
	return sender
}

// Start the main routine of the sender
//
// It sets the nonce and by the same occasion make sure that the remote node is accessible
func (s *Sender) Start() {
	go func(){
		<- s.ctx.Done()
		s.cancel() // Make sure the cancel function is triggered
	}()

	err := s.RefreshNonce()
	if err != nil {
		s.cancel() // Will close the sender as the client is not reachable
	}
}

// RefreshNonce gets the next nonce to use from the client
//
// It returns an error if the remote client could not be reached
func (s *Sender) RefreshNonce() error {
	nonce, err := s.client.PendingNonceAt(s.ctx, s.address)
	if err != nil {
		return err
	}
	s.nonce = nonce
	return nil
}

// SetSigner fetches the chainID from the client and cache it in the object
//
// It is necessary in order to 
func (s *Sender) SetSigner(signerT signerType, chainID *big.Int) error {
	s.chainID = chainID

	switch signerT {
	case EIP155:
		s.signer = types.NewEIP155Signer(chainID)
	case Homestead:
		s.signer = types.HomesteadSigner{}
	}

	return nil
}

// SuggestGasPrice returns an estimation of the gas price
//
// Is a simple wrapper around geth's SuggestGasPrice
func (s *Sender) SuggestGasPrice() (*big.Int, error) {
	return s.client.SuggestGasPrice(s.ctx)
}

// SendTransactionData publish a simple transaction to ethereum
//
// The function completely hides the responsibility to set
// the nonce, the gas price, the gasLimit and to sign the transaction
// to the user
func (s *Sender) SendTransactionData(to string, data []byte) error {
	return nil
}

// Sign a transaction in place
func (s *Sender) Sign(tx *types.Transaction) error {
	_, err := types.SignTx(tx, s.signer, &s.priv)
	if err!= nil {
		return err
	}
	return nil
}

// SendTransaction sends a transaction on the blockchain
func (s *Sender) SendTransaction(tx *Transaction) {

}
