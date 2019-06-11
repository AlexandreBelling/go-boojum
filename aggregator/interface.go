package aggregator

// Aggregator is a generic interface for aggregation services
type Aggregator interface {
	MakeExample() []byte
	AggregateTrees(left, right []byte) []byte
	Verify(buff []byte) bool
}
