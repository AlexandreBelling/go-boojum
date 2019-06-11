package aggregator

import "time"

// MockAggregator ..
type MockAggregator struct{}

// MakeExample ..
func (mock *MockAggregator) MakeExample() []byte {
	return []byte{0, 1, 2, 3, 4, 5, 6}
}

// AggregateTrees ..
func (mock *MockAggregator) AggregateTrees(left, right []byte) []byte {
	wakeup := time.After(time.Duration(1) * time.Second)
	<-wakeup
	return []byte{0, 1, 2, 3, 4, 5, 6}
}

// Verify ..
func (mock *MockAggregator) Verify(buff []byte) bool {
	return true
}
