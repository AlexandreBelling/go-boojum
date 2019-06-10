package election

import(
	"time"
	"github.com/golang/protobuf/proto"

	"github.com/AlexandreBelling/go-boojum/protocol"
	msg "github.com/AlexandreBelling/go-boojum/protocol/election/messages"
)

// Job contains the 
type Job struct{
	InputProofs		[][]byte
	Label			int
}

// Encode returns a marshalled job
func (j *Job) Encode() []byte {

	pb := &msg.AggregationJob{
		Type: "AggregationJob",        
		SubTrees: j.InputProofs,
		Label: int64(j.Label),
	}

	marshalled, err := proto.Marshal(pb)
	if err != nil {
		panic(err)
	} 
	return marshalled
}

// MarshalledJob is an alias for encoded proposal 
type MarshalledJob []byte


// Decode returns an unmarshalled Job
func (m MarshalledJob) Decode() (*Job, error) {
	j := &msg.AggregationJob{}
	err := proto.Unmarshal(m, j)
	if err != nil {
		return nil, err
	}

	return &Job{
		Label: int(j.Label),
		InputProofs: j.SubTrees,
	}, nil
}

// Proposal contains an data relative to an aggregation proposal
type Proposal struct {
	ID			protocol.ID
	Deadline	time.Time
}

// MarshalledProposal is an alias for encoded proposal 
type MarshalledProposal []byte

// Encode return a marshalled proposal
func (p *Proposal) Encode() []byte {

	pb := &msg.AggregationProposal{
		Id:			p.ID.String(),
		Deadline:	&msg.Timestamp{
			Sec: 		p.Deadline.Unix(),
			Nsec:		p.Deadline.UnixNano(),
		},
	}

	marshalled, err := proto.Marshal(pb)
	if err != nil {
		panic(err)
	} 
	return marshalled
}

// Decode attemps at returning an unmarshalled proposal
func (m MarshalledProposal) Decode() (*Proposal, error) {
	p := &msg.AggregationProposal{}
	err := proto.Unmarshal(m, p)
	if err != nil {
		return nil, err
	}

	return &Proposal{
		ID: 		protocol.StringToID(p.Id),
		Deadline: 	time.Unix(
			p.GetDeadline().Sec,
			p.GetDeadline().Nsec,
		),
	}, nil
}

// Result contains data relative to a result
type Result struct{
	Result			[]byte
	Label			int
	ID				protocol.ID
}

// Encode returns a marshalled Result
func (r *Result) Encode() []byte {

	pb := &msg.AggregationResult{
		Type: "AggregationResult",        
		Result: 	r.Result,
		Label: 		int64(r.Label),
		Id:			r.ID.String(),
	}

	marshalled, err := proto.Marshal(pb)
	if err != nil {
		panic(err)
	} 
	return marshalled
}

// MarshalledResult is an alias for encoded proposal 
type MarshalledResult []byte

// Decode returns an unmarshalled Result
func (m MarshalledResult) Decode() (*Result, error) {
	a := &msg.AggregationResult{}
	err := proto.Unmarshal(m, a)
	if err != nil {
		return nil, err
	}

	return &Result{
		Result: a.Result,
		Label: int(a.Label),
		ID: protocol.StringToID(a.Id),
	}, nil
}