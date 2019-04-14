package election

import (
	"github.com/golang/protobuf/proto"
	"github.com/AlexandreBelling/go-boojum/aggregator"
	msg "github.com/AlexandreBelling/go-boojum/protocol/election/messages"
)

// Worker ..
type Worker struct {
	Participant 		*Participant
	Tasks				[][]byte
	JobsIn				chan msg.AggregationRequest
	aggregator			aggregator.Aggregator
	LeaderAddress		string
}

// Run is the main routine of the worker
func (w *Worker) Run() {

	newBatch := <- w.Participant.Blockchain.NewBatch
	w.IdentifyLeader(newBatch)

	for {

		select {

		case job := <- w.JobsIn:
			result := w.DoJob(job.GetSubTrees())
			w.SendResult(result, job.GetToken())

		case <- w.Participant.Blockchain.BatchDone:
			return
		}
	}
}

// SendProposal ..
func (w *Worker) SendProposal() error {

	msg := &msg.AggregationProposal{
			Address: w.Participant.Address,
			Signature: []byte{},
		}

	marshalled, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}

	w.Participant.Network.Send(marshalled, w.LeaderAddress)
	return nil
}

// IdentifyLeader ..
func (w *Worker) IdentifyLeader(newBatch [][]byte) {
	w.LeaderAddress = "0"
}

// DoJob ..
func (w *Worker) DoJob(job [][]byte) []byte {
	return w.aggregator.AggregateTrees(job[0], job[1])
}

// SendResult ..
func (w* Worker) SendResult(result []byte, token int64) error {

	msg := &msg.AggregationResult{
		Result: result,
		Token: token,
	}

	marshalled, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}

	w.Participant.Network.Send(marshalled, w.LeaderAddress)
	return nil
}