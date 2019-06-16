package ethereum

import (
	"math/big"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
)

// Transaction is a wrapper for geth transaction
type Transaction types.Transaction

// NewTransaction returns an ethereum transaction
func NewTransaction(nonce uint64, to string, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte ) *Transaction {
	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(to),
		amount, gasLimit, gasPrice, data,
	)
	return (*Transaction)(tx)
}

// GetTo returns the recipient of the transaction
func (t *Transaction) GetTo() string {
	return (*types.Transaction)(t).To().Hex()
}

// GetData returns the recipient of the transaction
func (t *Transaction) GetData() []byte {
	return (*types.Transaction)(t).Data()
}
