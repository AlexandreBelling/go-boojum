package election

import(
	"fmt"

	"github.com/AlexandreBelling/go-boojum/protocol"
	"github.com/AlexandreBelling/go-boojum/network"
)

var (
	// TopicPath is the basis of every topic used by this protocol
	TopicPath 				= "boojum.protocol.election"
	// ProposalTopicPath is the base string used for the proposals
	ProposalTopicPath 		string
	// ResultTopicPath is the base string used for the aggregation results
	ResultTopicPath 		string
	// JobTopicPath is the base string used for the aggregation requests
	JobTopicPath 			string
)

func init() {
	ProposalTopicPath		= fmt.Sprintf("%v.proposal", TopicPath)
	ResultTopicPath 		= fmt.Sprintf("%v.result", TopicPath)
	JobTopicPath 			= fmt.Sprintf("%v.request", TopicPath)
}

// JobTopic format a topic name for a given job and round
func JobTopic(roundID protocol.ID, workerID protocol.ID) string {
	return fmt.Sprintf("%v.%v.%v", 
		JobTopicPath, 
		roundID, 
		workerID,
	)
}

// ResultTopic format a topic name for a given result
func ResultTopic(roundID protocol.ID, label int) string {
	return fmt.Sprintf("%v.%v.%v",
		ResultTopicPath,
		roundID,
		label,
	)
}

// ProposalTopic format a topic name for a new proposal
func ProposalTopic(roundID protocol.ID) string {
	return fmt.Sprintf("%v.%v",
		ProposalTopicPath,
		roundID,
	)
}

// TopicProvider is a helper to interact with topics in a friendly way
type TopicProvider struct {
	Network 	network.PubSub
	Round		*Round
}

// ProposalTopic returns the appropriate result topic
func (t *TopicProvider) ProposalTopic() (network.Topic) {
	return t.Network.GetTopic(
		ProposalTopic(t.Round.ID),
	)
}

// JobTopic returns the appropriate job topic
func (t *TopicProvider) JobTopic(ID protocol.ID) (network.Topic) {
	return t.Network.GetTopic(
		JobTopic(t.Round.ID, ID),
	)
}

// ResultTopic in the appropriate result topic
func (t *TopicProvider) ResultTopic(label int) (network.Topic) {
	return t.Network.GetTopic(
		ResultTopic(t.Round.ID, label),
	)
}

// PublishProposal to alert the leader, we are ready
func (t *TopicProvider) PublishProposal(p *Proposal) error {
	return t.Network.Publish(
		ProposalTopic(t.Round.ID),
		p.Encode(),
	)
}

// PublishJob to make a worker do it
func (t *TopicProvider) PublishJob(j *Job, ID protocol.ID) error {
	return t.Network.Publish(
		JobTopic(t.Round.ID, ID),
		j.Encode(),
	)
}

// PublishResult to alert a leader that we are done
func (t *TopicProvider) PublishResult(r *Result) error {
	return t.Network.Publish(
		ResultTopic(t.Round.ID, r.Label),
		r.Encode(),
	)
}