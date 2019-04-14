package aggregator

// MockAggregator ..
type MockAggregator struct {}

// MakeExample ..
func (mock *MockAggregator) MakeExample() []byte {
	return []byte{0, 1, 2, 3, 4, 5, 6}
}

// AggregateTrees ..
func (mock *MockAggregator) AggregateTrees(left, right []byte) []byte {
	return []byte{0, 1, 2, 3, 4, 5, 6}
}

// Verify ..
func (mock *MockAggregator) Verify(buff []byte) bool {
	return true
}