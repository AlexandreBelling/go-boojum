package ethereum

import (
	"github.com/AlexandreBelling/go-boojum/blockchain/ethereum"
)

// Agent (ethereum.Agent) implements election.BlockchainAgent in the case of ethereum
type Agent struct {
	sender 		*ethereum.Sender
	listener	*ethereum.Listener
}

// NotifyOnNewBatch implements blockchain.Notification and 
// verifies if the transaction triggers a new batch and callback the Participant if so
func (e *EthereumAgent) NotifyOnNewBatch(blockchain.Transaction) {}

// NotifyOnBatchDone implements blockchain.Notification and
// verifies if a transaction indicate that an aggregation step has been completed
func (e *EthereumAgent) NotifyOnBatchDone(blockchain.Transaction) {}

