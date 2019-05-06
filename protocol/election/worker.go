package election

import (
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	msg "github.com/AlexandreBelling/go-boojum/protocol/election/messages"
)

// Worker ..
type Worker struct {
	Participant 		*Participant
	Tasks				[][]byte
	JobsIn				chan msg.AggregationRequest
	LeaderAddress		string
}

// Run is the main routine of the worker
func (w *Worker) Run() {

	w.IdentifyLeader(w.Tasks)
	w.SendProposal()

	for {
		select {

		case job := <- w.JobsIn:
			log.Debugf("Boojum | Worker %v | Acknowledging %v", w.Participant.Address, job.Token)
			result := w.DoJob(job.GetSubTrees())
			w.SendResult(result, job.GetToken())
			w.SendProposal()

		case <- w.Participant.Blockchain.BatchDone:
			log.Debugf("Boojum | Worker %v | Batch done", w.Participant.Address)
			return
		}
	}
}

// NewWorker construct a new worker object
func NewWorker(participant *Participant, tasks [][]byte) *Worker {

	return &Worker{
		Participant:	participant,
		Tasks:			tasks,
		LeaderAddress:	participant.Blockchain.GetLeaderAddress(tasks),

		JobsIn:			make(chan msg.AggregationRequest),
	}
}

// SendProposal ..
func (w *Worker) SendProposal() error {

	log.Debugf("Boojum | Worker: %v | Sending proposal", w.Participant.Address)

	msg := &msg.AggregationProposal{
			Type: "Proposal",
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
	return w.Participant.Aggregator.AggregateTrees(job[0], job[1])
}

// SendResult ..
func (w* Worker) SendResult(result []byte, token int64) error {

	msg := &msg.AggregationResult{
		Type: "Result",
		Result: result,
		Token: token,
	}

	marshalled, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}

	w.Participant.Network.Send(marshalled, w.LeaderAddress)
	log.Debugf("Boojum | Worker : %v | Sending result for %v", w.Participant.Address, msg.Token)
	return nil
}