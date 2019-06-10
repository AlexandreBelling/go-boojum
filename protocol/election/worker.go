package election

import(
	"context"
)

// Worker contains all the logic required to aggregate the proofs
type Worker struct {
	ctx			context.Context
	cancel		context.CancelFunc
	Round		*Round
}

// NewWorker returns a newly constructed worker
func NewWorker(r *Round) *Worker {
	return &Worker{
		ctx: 		r.ctx,
		cancel:		r.cancel,
		Round:		r,
	}
}

// Aggregate performs an aggregation
func (w *Worker) Aggregate(job *Job) (*Result, error) {

	data := w.Round.Participant.Aggregator.AggregateTrees(
		job.InputProofs[0], job.InputProofs[1], // TODO: Support multi-arity	
	)

	return &Result{
		Result: data,
		Label: job.Label,
		ID:	w.Round.Participant.ID,
	}, nil
}

// PublishProposal to alert the leader, we are ready
func (w *Worker) PublishProposal() error {
	proposal := &Proposal{ ID: w.Round.Participant.ID }
	return w.Round.TopicProvider.PublishProposal(proposal)
}

// Start the Worker routine
func (w *Worker) Start() error {
	topic := w.Round.TopicProvider.JobTopic(
		w.Round.Participant.ID,
	)

	jobChan, err := topic.Chan()
	if err != nil {
		return err
	}
	
	go func(){
		defer w.cancel()

		for {
			err := w.PublishProposal()
			if err != nil {
				return
			}
	
			select {
			case <- w.ctx.Done():
				return
	
			case jobEncoded := <- jobChan:
				job, err := MarshalledJob(jobEncoded).Decode()
				if err != nil {
					return
				}
	
				res, err := w.Aggregate(job)
				if err != nil {
					return
				}

				_ = w.Round.TopicProvider.PublishResult(res)
			}
		}
	}()

	return nil
}