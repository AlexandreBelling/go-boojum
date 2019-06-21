package election

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// Task is the function processing the job
type Task func(context.Context, *Proposal) error

// JobPool schedules the jobs
type JobPool struct {
	ctx           context.Context
	proposalQueue chan *Proposal
}

// NewJobPool returns a job pool object
func NewJobPool(ctx context.Context) *JobPool {
	return &JobPool{
		ctx:           ctx,
		proposalQueue: make(chan *Proposal, 1024),
	}
}

// EnqueueProposal enqueue a proposal in the channel
func (j *JobPool) EnqueueProposal(ctx context.Context, p *Proposal) error {
	select {
	case j.proposalQueue <- p:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// DequeueProposal returns a proposal waits for a proposal to be dequeued
func (j *JobPool) DequeueProposal(ctx context.Context) (*Proposal, error) {
	for {
		select {
		case p := <-j.proposalQueue:
			if p.Deadline.Before(time.Now()) {
				log.Infof("Got an expired proposal. The deadline was at %v and it is %v", p.Deadline, time.Now())
				continue // Skip outdated proposal
			}
			return p, nil
		case <-ctx.Done():
			return &Proposal{}, fmt.Errorf("Could not dequeue proposal, proposal queue is empty until context expired")
		}
	}
}

// AddJob adds the job to the pool until it is completed.
// It assumes the job is well-formed.
func (j *JobPool) AddJob(jobctx context.Context, task Task) {
	go j.addJobSync(jobctx, task)
}

// addJobSync handle a new job in a synchronous way
func (j *JobPool) addJobSync(jobctx context.Context, task Task) {
	for {
		// Wait without timeout for a new proposal to arrive
		select {
		default:
			propal, _ := j.DequeueProposal(j.ctx)
			log.Infof("Got proposal from : %v", propal.ID.Pretty())
			err := task(jobctx, propal)
			if err == nil {
				return
			}
		case <-j.ctx.Done():
			return
		}
	}
}
