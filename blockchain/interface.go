package blockchain

import(
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// Client is an interface to Ethereum
type Client interface{
	SendTransactionRetry(tx *ethtypes.Transaction)
}