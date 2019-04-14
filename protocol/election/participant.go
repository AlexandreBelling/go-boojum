package election

import (
	"github.com/AlexandreBelling/go-boojum/aggregator"
	net "github.com/AlexandreBelling/go-boojum/network"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/AlexandreBelling/go-boojum/blockchain"
	msg "github.com/AlexandreBelling/go-boojum/protocol/election/messages"
	"github.com/golang/protobuf/proto"
)

// Participant ..
type Participant struct {
	Address 				string
	ParticipantAddresses	[]string
	
	Network 				net.Network
	Blockchain				*ParticipantBlockchainInterface
	Aggregator				aggregator.Aggregator
	
	NetworkIn 				<-chan []byte
	StopChan				chan bool
}

// NewParticipant ..
func NewParticipant(address string) *Participant {
	return &Participant{Address: address}
}

// WithNetwork ..
func (par *Participant) WithNetwork(network net.Network) *Participant {
	par.Network = network
	par.NetworkIn = network.GetConnexion(par.Address)
	return par
}

// WithBCInterface ..
func (par *Participant) WithBCInterface(blockchain *ParticipantBlockchainInterface) (*Participant) {
	par.Blockchain = blockchain
	return par
}

// WithAggregator ..
func (par *Participant) WithAggregator(aggregator aggregator.Aggregator) (*Participant) {
	par.Aggregator = aggregator
	return par
}

// GetRank returns the rank of the participant in the given round
// This will be usefull to support crash failure
func (par *Participant) GetRank(batch [][]byte) int {

	if par.Blockchain.GetLeaderAddress(batch) == par.Address {
		return 0
	}
	return 1
}

// Run starts the main routine of the participant
func (par *Participant) Run() {
	
	// Main routine
	go func() {
		for {
			batch := <- par.Blockchain.NewBatch
			switch par.GetRank(batch) {
			default:
				// Do workers stuff
				// Blocking call
			case 0:
				// Do leader stuff
			}
		}
	}()

	// Receiving on this chanel triggers early exit
	<- par.StopChan
	return
}

// ForwardNetworkToWorker ..
func (par *Participant) ForwardNetworkToWorker(worker *Worker, quit <-chan bool) {
	for {
		select{

		case <-quit:
			return

		case marshalled := <-par.NetworkIn:	
			request := &msg.AggregationRequest{}
			err := proto.Unmarshall(marshalled, request)
			if err == nil {
				worker.JobsIn <- request
			}

		}

	}
}

// ForwardNetworkToWorker ..
func (par *Participant) ForwardNetworkToLeader(leader *Leader, quit <-chan bool) {
	for {
		select{

		case <-quit:
			return

		case marshalled := <-par.NetworkIn:

			proposal := &msg.AggregationProposal{}
			err := proto.Unmarshall(marshalled, request)
			if err == nil {
				leader.ProposalsChan <- request
				continue
			}

			result := &msg.AggregationResult{}
			err = proto.Unmarshall(marshalled, request)
			if err == nil {
				leader.ResultsChan <- result
				continue
			}
		}
		
	}
}

// ParticipantBlockchainInterface is a component of Participant. 
// It manages its communications with the blockchain.
type ParticipantBlockchainInterface struct {

	BlockStream			chan ethtypes.Block
	NewBatch 			chan [][]byte
	BatchDone 			chan bool

	Backend 			blockchain.Client
	ContractAddress		ethcommon.Address
}

// PublishAggregated ..
func (bci *ParticipantBlockchainInterface) PublishAggregated(aggregated []byte) {
	// Do stuffs with to make sure the transaction is eventually mined
	bci.Backend.SendTransactionRetry(&ethtypes.Transaction{})
	return
}

// GetLeaderAddress returns the address of the leader for the current batch
func (bci *ParticipantBlockchainInterface) GetLeaderAddress(tasks [][]byte) string {
	return "0"
}