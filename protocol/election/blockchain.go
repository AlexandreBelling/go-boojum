package election

import (
	"time"
	log "github.com/sirupsen/logrus"
	"github.com/AlexandreBelling/go-boojum/blockchain"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// ParticipantBlockchainInterface is a component of Participant. 
// It manages its communications with the blockchain.
type ParticipantBlockchainInterface struct {

	// Todo pass transactions instead of entires blocks
	BlockStream			chan ethtypes.Block
	NewBatch 			chan [][]byte
	BatchDone 			chan bool

	Backend 			blockchain.Client
	ContractAddress		ethcommon.Address
}

// NewParticipantBlockchainInterface construct a new ParticipantBlockchainInterface
func NewParticipantBlockchainInterface(blockchain blockchain.Client) *ParticipantBlockchainInterface {
	return &ParticipantBlockchainInterface{
		Backend: 		blockchain,

		BlockStream: 	make(chan ethtypes.Block, 32),
		NewBatch: 		make(chan [][]byte, 32),
		BatchDone: 		make(chan bool, 32),
	}
}

// PublishAggregated ..
func (bci *ParticipantBlockchainInterface) PublishAggregated(aggregated []byte) {
	// Do stuffs with to make sure the transaction is eventually mined
	log.Infof("Boojum | Worker | Publish on-chain | Time : %v", time.Now())
	bci.Backend.SendTransactionRetry(&ethtypes.Transaction{})
	return
}

// GetLeaderAddress returns the address of the leader for the current batch
func (bci *ParticipantBlockchainInterface) GetLeaderAddress(tasks [][]byte) string {
	return "0"
}