package layered

import (
	msg "github.com/AlexandreBelling/go-boojum/protocol/layered/messages"
	"github.com/golang/protobuf/proto"
	"math"
	"math/big"
	"crypto/rand"
	"fmt"
	"time"
)


// Round ..
type Round struct {}

// Role ..
type Role interface {}

// ResultDemuxer maps token with the right result channel
type ResultDemuxer struct {
	ResultsPerID 	map[int64]chan []byte
}

// NewResultDemuxer is a constructor for ResultDemuxer
func NewResultDemuxer() (*ResultDemuxer) {
	return &ResultDemuxer{
		ResultsPerID: make(map[int64]chan []byte),
	}
}

// AddConsumer adds a channel in the demultiplexer map
func (demux *ResultDemuxer) AddConsumer(token int64, consumingChannel chan []byte) {
	demux.ResultsPerID[token] = consumingChannel
	return
}

// FanOut receive a message and pass it to the right channel
func (demux *ResultDemuxer) FanOut(token int64, data []byte) error {

	if demux.ResultsPerID[token] == nil {
		return fmt.Errorf("Channel for token %v does not exists", token)
	}

	demux.ResultsPerID[token] <- data
	return nil
}

// Leader ..
type Leader struct {

	Timeout int
	Arity int

	ProposalsChan 		chan msg.AggregationProposal
	ResultsChan			chan msg.AggregationResult
	
	ResultMux  			*ResultDemuxer
	Root				*Tree
	Participant 		*Participant
}

// NewLeader ..
func NewLeader(tasks [][]byte, arity int, nWorkers int, participant *Participant) (l *Leader) {

	height := int(
		math.Log2(
			float64(
				len(tasks),		
	))) + 1

	root := NewTree(height, arity)
	leaves := root.GetLeaves()
	for index := range leaves {
		leaves[index].payloadChan <- tasks[index]
	}

	l = &Leader{
		Timeout: 10,
		Arity: 2,

		ProposalsChan: 	make(chan msg.AggregationProposal, 	nWorkers),
		ResultsChan:	make(chan msg.AggregationResult, 	height),

		ResultMux: 		NewResultDemuxer(),
		Root:			NewTree(height, arity),
		Participant:	participant,		
	}
	
	return l
}

// Schedule waits for operands to be aggregated then add to scheduler
func (l *Leader) Schedule(t *Tree) {

	// For non-leaf nodes only
	if len(t.children) > 0 {

		// Recurse the call in children
		for _, child := range l.Root.children {
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
		l.DispatchRetry(t, payloads)
	}

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

		timeoutChan := make(chan bool)
		go func() {
			time.Sleep(time.Duration(l.Timeout) * time.Second)
			timeoutChan <- true
		}()

		// Wait for the tasks to be completed or timeout
		for {
			select {
			case result:= <- doneChan:
				t.payloadChan <- result
				// TODO: If verification is light enough to be heavily multithreaded, add it here
				return
			case <- timeoutChan:
				continue taskLoop
			}
		}

	}
}

// Dispatch dispatch a job to worker
func (l *Leader) Dispatch(token int64, payloads [][]byte, address string) error {

	request := msg.AggregationRequest{
		SubTrees: payloads,
		Token: token,
	}
	marshalled, err := proto.Marshal(&request)
	if err != nil {
		return err
	}
	return l.Participant.Network.Send(marshalled, address)
}

// DistributeResults ..
func (l *Leader) DistributeResults() {
	for {
		result := <- l.ResultsChan
		// TODO: If verification so heavy, that it needs its own thread policy add it here
		l.ResultMux.FanOut(result.GetToken(), result.GetResult())
	}
}

// RandomToken is panicable function that returns an error
func RandomToken() int64 {
	
	nBig, err := rand.Int(rand.Reader, big.NewInt(64))
	if err != nil {
		panic(err)
	}
	return nBig.Int64()
}