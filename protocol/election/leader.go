package election

import (
	msg "github.com/AlexandreBelling/go-boojum/protocol/election/messages"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"math"
	"math/big"
	"crypto/rand"
	"fmt"
	"time"
	"sync"
)

// Leader ..
type Leader struct {

	Participant 	*Participant
	Tasks			[][]byte
	Timeout 		int
	Arity 			int

	ProposalsChan 	chan msg.AggregationProposal
	ResultsChan		chan msg.AggregationResult
	
	ResultMux  		*ResultDemuxer
	Root			*Tree

	Stop			chan bool
}

// NewLeader construct a new leader instance
func NewLeader(tasks [][]byte, arity int, nWorkers int, participant *Participant) (l *Leader) {

	height := int(
		math.Log2(
			float64(
				len(tasks),		
	))) + 1

	l = &Leader{

		Participant:	participant,
		Tasks:			tasks,
		Timeout: 		1,
		Arity: 			2,

		ProposalsChan: 	make(chan msg.AggregationProposal, 	nWorkers),
		ResultsChan:	make(chan msg.AggregationResult, 	len(tasks)),

		ResultMux: 		NewResultDemuxer(),
		Root:			NewTree(height, arity),
		
		Stop: 			make(chan bool),
	}
	
	return l
}

// Run contains the main routine of a Leader instance
func (l *Leader) Run() {

	// Start the multiplexer goroutine
	stopDemux := make(chan bool)
	go l.DemuxResults(stopDemux)
	defer func(){ stopDemux <- true }()

	go l.Schedule(l.Root)
	l.Start()  // This triggers the leader
	
	// When it's finished publish it on the chain
	aggregated := <- l.Root.payloadChan
	l.Participant.Blockchain.PublishAggregated(aggregated)
	<- l.Participant.Blockchain.BatchDone
	return
}

// Start dispatching jobs to the workers
func (l *Leader) Start() {
	// Starts
	leaves := l.Root.GetLeaves()
	for index := range leaves {
		leaves[index].payloadChan <- l.Tasks[index]
	}
}

// Schedule waits for operands to be aggregated then add to scheduler
func (l *Leader) Schedule(t *Tree) {

	// For leaf nodes only
	if len(t.children) == 0 {
		// Wait to be assigned a payload
		t.payload = <- t.payloadChan
		return 
	}

	// Recurse the call in children
	for _, child := range t.children {
		if child != nil {
			go l.Schedule(child)
		}
	}
	
	// Block until all subtasks have been completed. 
	payloads := make([][]byte, len(t.children))
	for index, child := range t.children {
		if child != nil {
			// We need to assign in a separate variable 
			// first before passing to the child attribute
			payload := <- child.payloadChan
			payloads[index] = payload
		}
	}

	// Ensures the task endup being done
	log.Debugf("Boojum | Leader | Handling a new tasks")
	l.DispatchRetry(t, payloads)
	return 
}

// DispatchRetry make sure the job is completed
func (l *Leader) DispatchRetry(t *Tree, payloads [][]byte) {

	token := RandomToken()
	doneChan := make(chan []byte)
	l.ResultMux.AddConsumer(token, doneChan)

	taskLoop:
	for {

		proposal:= <- l.ProposalsChan
		err := l.Dispatch(token, payloads, proposal.Address)
		if err != nil {
			continue taskLoop
		}

		timer := time.NewTimer(time.Duration(l.Timeout) * time.Second)
		defer timer.Stop()

		// Wait for the tasks to be completed or timeout
		for {
			select {
			case result:= <- doneChan:
				l.Participant.Aggregator.Verify(result)
				// TODO: Add a way to test wether the proof is the right proof
				t.payloadChan <- result
				return
			case <- timer.C:
				log.Debugf("Boojum | Leader | Got a timeout for %v", token)
				continue taskLoop
			}
		}

	}
}

// Dispatch dispatch a job to worker
func (l *Leader) Dispatch(token int64, payloads [][]byte, address string) error {

	request := msg.AggregationRequest{
		Type: "Request",
		SubTrees: payloads,
		Token: token,
	}

	marshalled, err := proto.Marshal(&request)
	log.Debugf("Boojum | Leader | Dispatching %v", token)
	if err != nil {
		return err
	}

	return l.Participant.Network.Send(marshalled, address)
}

// DemuxResults is an auxilliary routine
func (l *Leader) DemuxResults(quit chan bool) {
	for {
		select {

		case <- quit:
			close(quit)
			return

		case result := <- l.ResultsChan:
			l.ResultMux.FanOut(result.GetToken(), result.GetResult())
		}
	}
}

// RandomToken is panicable function that returns an error
func RandomToken() int64 {
	
	nBig, err := rand.Int(rand.Reader, big.NewInt(1<<32))
	if err != nil {
		panic(err)
	}
	return nBig.Int64()
}

// ResultDemuxer maps token with the right result channel
type ResultDemuxer struct {
	ResultsPerID 	map[int64]chan []byte
	mut 			sync.Mutex
}

// NewResultDemuxer is a constructor for ResultDemuxer
func NewResultDemuxer() (*ResultDemuxer) {
	return &ResultDemuxer{
		ResultsPerID: make(map[int64]chan []byte),
	}
}

// AddConsumer adds a channel in the demultiplexer map
func (demux *ResultDemuxer) AddConsumer(token int64, consumingChannel chan []byte) {
	demux.mut.Lock()
	demux.ResultsPerID[token] = consumingChannel
	demux.mut.Unlock()
	return
}

// FanOut receive a message and pass it to the right channel
func (demux *ResultDemuxer) FanOut(token int64, data []byte) error {
	demux.mut.Lock()
	if demux.ResultsPerID[token] == nil {
		return fmt.Errorf("Channel for token %v does not exists", token)
	}

	demux.ResultsPerID[token] <- data
	demux.mut.Unlock()
	return nil
}