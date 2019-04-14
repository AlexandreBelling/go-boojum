package election

import (
	log "github.com/sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/AlexandreBelling/go-boojum/aggregator"
	"github.com/AlexandreBelling/go-boojum/blockchain"
	net "github.com/AlexandreBelling/go-boojum/network"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	msg "github.com/AlexandreBelling/go-boojum/protocol/election/messages"
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
func (par *Participant) WithBCInterface(blockchain blockchain.Client) (*Participant) {

	par.Blockchain = &ParticipantBlockchainInterface{
		Backend: 		blockchain,

		BlockStream: 	make(chan ethtypes.Block),
		NewBatch: 		make(chan [][]byte),
		BatchDone: 		make(chan bool),
	}

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
			println("New batch")
			switch par.GetRank(batch) {
			default:
				log.Infof("Boojum | Participant: %v | Started worker", par.Address)
				w := NewWorker(par, batch)
				w.Run()
			case 0:
				log.Infof("Boojum | Participant: %v | Started leader", par.Address)
				l := NewLeader(batch, 2, 10, par)
				l.Run()
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
			err := proto.Unmarshal(marshalled, request)
			if err == nil {
				worker.JobsIn <- *request
			}
		}

	}
}

// ForwardNetworkToLeader ..
func (par *Participant) ForwardNetworkToLeader(leader *Leader, quit <-chan bool) {
	
	for {
		select{

		case <-quit:
			return

		case marshalled := <-par.NetworkIn:

			proposal := &msg.AggregationProposal{}
			err := proto.Unmarshal(marshalled, proposal)
			if err == nil {
				leader.ProposalsChan <- *proposal
				continue
			}

			result := &msg.AggregationResult{}
			err = proto.Unmarshal(marshalled, result)
			if err == nil {
				leader.ResultsChan <- *result
				continue
			}
		}
		
	}
}

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