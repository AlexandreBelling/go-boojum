package election

import(
	"fmt"
)

var (
	// TopicPath is the basis of every topic used by this protocol
	TopicPath 			= "boojum.protocol.election"
	// ProposalTopic is the topic used for the proposals
	ProposalTopic 		string
	// ResultTopicPath is the base string used for the aggregation results
	ResultTopicPath 	string
	// RequestTopicPath is the base string used for the aggregation requests
	RequestTopicPath 	string
)

func init() {
	ProposalTopic 		= fmt.Sprintf("%v.proposal", TopicPath)
	ResultTopicPath 	= fmt.Sprintf("%v.result", TopicPath)
	RequestTopicPath 	= fmt.Sprintf("%v.request", TopicPath)
}