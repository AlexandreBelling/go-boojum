package election

import(
	"fmt"
	"time"
	"context"

	"github.com/AlexandreBelling/go-boojum/protocol"
)

// Leader is a control struct that schedule the aggregation process
type Leader struct {

	JobPool 		*JobPool
	Tree			*protocol.Tree // Root of the aggregation tree
	Round			*Round
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
			l.JobPool.AddJob(context.Background(), task)
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

	handler := func(_ context.Context, p Proposal) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1) * time.Minute)
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

// Start the leader routine
func (l *Leader) Start() {}
