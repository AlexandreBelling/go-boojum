package election

import(
	"fmt"
	"time"
	"math"
	"context"

	// log "github.com/sirupsen/logrus"

	"github.com/AlexandreBelling/go-boojum/protocol"
)

// Leader is a control struct that schedule the aggregation process
type Leader struct {
	ctx				context.Context

	JobPool 		*JobPool
	Tree			*protocol.Tree // Root of the aggregation tree
	Round			*Round

	cancel			context.CancelFunc			
}

// Start the leader routine
func (l *Leader) Start() {
	// By populating the leaves of the tree, the aggregation scheduling process is triggered through the
	
	l.populateLeaves()
	l.ListenForProposal()
}

// NewLeader constructs a new leader
func NewLeader(r *Round) *Leader {
	l := &Leader{
		ctx: 		r.ctx,
		cancel:		r.cancel,
		JobPool: 	NewJobPool(r.ctx), 
		Round: 		r, 
	}

	l.Tree = protocol.NewTree(l.getTreeHeight(), 2)
	InitializeNodes(l.Tree, 
		l.OnReadinessUpdateHook(),
		l.OnRootProofUpdateHook(),
	)

	return l
}

// OnRootProofUpdateHook returns a hook that is triggered when setting a value to the root node's aggregated proof
// Taking action to publish it on-chain
func (l *Leader) OnRootProofUpdateHook() NodeHook {
	return func(n *Node) {
		l.PublishOnChain(n.AggregateProof)
	}
}

// OnReadinessUpdateHook returns a HookOnReadinessUpdate 
// that sends a new job to the pool if ready
// It is automaticallt triggered by Node when all its children are completed 
func (l *Leader) OnReadinessUpdateHook() NodeHook {
	return func(n *Node) {

		task, err := l.MakeJobHandler(n)
		if err != nil {
			panic(err)
		}

		if n.IsReady() {
			l.JobPool.AddJob(l.ctx, task)
		}

		// log.Infof("Pushed job to the pool")
	}
}

// MakeJobHandler returns a jobpool tasks that handle an aggregation job 
func (l *Leader) MakeJobHandler(n *Node) (Task, error) {

	if n.Topic != nil { // TODO: Remove when are sure, the code is correct
		panic(fmt.Errorf("Attempted to create two topic for the same node"))
	}
	
	// No way we can get the error here since we just instantiated the topic
	rtopic := l.Round.TopicProvider.ResultTopic(l.ctx, n.Label)
	defer rtopic.Close()

	resultChan, err := rtopic.Chan()
	if err != nil {
		return nil, err
	}

	handler := func(ctx context.Context, p *Proposal) error {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(1) * time.Minute)
		defer cancel()

		err = l.Round.TopicProvider.PublishJob(n.Job(), p.ID)
		if err != nil {
			return err
		}

		select {

		case <- ctx.Done():
			return ctx.Err()

		case r := <- resultChan:
			result, err := MarshalledResult(r).Decode()
			if err != nil {
				return err
			}

			n.SetAggregateProof(result.Result)
			// log.Infof("Got an aggregated result")
			return nil
		}
	}
	
	return handler, nil
}

func (l *Leader) getTreeHeight() int {
	// At the moment, the arity size is stuck to two
	batchSize := len(l.Round.Batch)
	if batchSize == 0 {
		return 1
	}

	return int(
		math.Ceil(math.Log2(
			float64(batchSize),
		)),
	) + 1
}

func(l *Leader) populateLeaves() {
	leaves := l.Tree.GetLeaves()
	batch := l.Round.Batch

	for i:=0; i<len(batch); i++ {
		leaves[i].Node.(*Node).SetAggregateProof(batch[i])
	}

	if len(batch) == len(leaves) { 
		return 
	}

	// Pad the remaining leaves with dummy proofs
	examples := l.Round.Participant.Aggregator.MakeExample()
	for i:=len(batch); i<len(leaves); i++ {
		leaves[i].Node.(*Node).SetAggregateProof(examples)
	}
}

// ListenForProposal start an async loop fetching new proposal and enqueuing them
func(l *Leader) ListenForProposal() error {
	
	topic := l.Round.TopicProvider.ProposalTopic(l.ctx)
	topicChan, err := topic.Chan()
	if err != nil {
		return err
	}

	go func(){
		defer topic.Close()
		for {
			select {
			case <- l.ctx.Done():
				return

			case b, ok := <- topicChan:
				if !ok { return } // If the topic was closed for another reason
				decoded, err := MarshalledProposal(b).Decode()
				if err != nil { continue }
				l.JobPool.EnqueueProposal(l.ctx, decoded)
				// log.Infof("Got a new proposal")
			}
		}
	}()

	return nil
}

// PublishOnChain sends the aggregated proof on-chain
func (l *Leader) PublishOnChain(aggregatedProof []byte) {
	l.Round.Participant.Blockchain.PublishAggregated(aggregatedProof)
	// log.Infof("Published the result onchain")
}