package election

import(
	"github.com/AlexandreBelling/go-boojum/aggregator"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	//log "github.com/sirupsen/logrus"
)

// BCClientMock is a mock for a blockchain
type BCClientMock struct {
	Participants 	[]BCUser
	Batch 	 		[][]byte
	BatchSize		int
}

// GenerateBatch creates a batch of zksnark proofs to be aggregated
// They are all identicals, this is a batch of 32 proofs
func (mock *BCClientMock) GenerateBatch(boojum aggregator.Aggregator) {
	leaf := boojum.MakeExample()
	mock.Batch = make([][]byte, 8)
	
	for i:=0; i<2; i++ {
		mock.Batch[i] = leaf
	}
}

// Connect creates a transaction with the blockchain interface
func (mock *BCClientMock) Connect(blockchainInterface *BCUser) {
	mock.Participants = append(mock.Participants, *blockchainInterface)
	return
}

// SendTransactionRetry will ensure the transaction passes. It is a Mock
// tx argument is never read and is supposed to be sent to nil
func (mock BCClientMock) SendTransactionRetry(tx *ethtypes.Transaction) {
	for _, participant := range mock.Participants {
		participant.BatchDone 	<- true
		participant.NewBatch 	<- mock.Batch 
	}
}

// NewBatch initiate the test
func (mock *BCClientMock) NewBatch() {
	for _, participant := range mock.Participants {
		participant.NewBatch 	<- mock.Batch 
	}
}