package election

import(
	"fmt"
	"time"
	"math"
	"context"

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
}

// NewLeader constructs a new leader
func NewLeader(ctx context.Context, r *Round) *Leader {
	ctx, cancel := context.WithCancel(ctx)

	l := &Leader{
		ctx: 		ctx,
		cancel:		cancel,
		JobPool: 	NewJobPool(ctx), 
		Round: 		r, 
	}

	l.Tree = protocol.NewTree(l.getTreeHeight(), 2)
	InitializeNodes(l.Tree, l.OnreadinessUpdateHook())

	return l
}

// OnreadinessUpdateHook returns a HookOnReadinessUpdate 
// that sends a new job to the pool if ready 
func (l *Leader) OnreadinessUpdateHook() HookOnReadinessUpdate {
	return func(n *Node) {

		task, err := l.MakeJobHandler(n)
		if err != nil {
			panic(err)
		}

		if n.IsReady() {
			l.JobPool.AddJob(l.ctx, task)
		}
	}
}

// MakeJobHandler returns a jobpool tasks that handle an aggregation job 
func (l *Leader) MakeJobHandler(n *Node) (Task, error) {

	if n.Topic != nil { // TODO: Remove when are sure, the code is correct
		panic(fmt.Errorf("Attempted to create two topic for the same node"))
	}

	topicResult, err := l.Round.Participant.Network.GetTopic(
		fmt.Sprintf("%v.%v", ResultTopicPath, n.Label),
	)
	n.Topic = topicResult

	if err != nil {
		return nil, err
	}
	
	// No way we can get the error here since we just instantiated the topic
	resultChan, _ := topicResult.Chan()
	jobEncoded := n.Job().Encode()

	handler := func(ctx context.Context, p *Proposal) error {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(1) * time.Minute)
		defer cancel()

		err = l.Round.Participant.Network.Publish(
			fmt.Sprintf("%v.%v", RequestTopicPath, p.ID),
			jobEncoded,
		)

		select {

		case <- ctx.Done():
			return ctx.Err()

		case r := <- resultChan:
			result, err := MarshalledResult(r).Decode()
			if err != nil {
				return err
			}

			n.SetAggregateProof(result.Result)
			topicResult.Close()
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
	
	topic, err := l.Round.Participant.Network.GetTopic(ProposalTopic)
	if err != nil {
		return err
	}

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
			}
		}
	}()

	return nil

}
