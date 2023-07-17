package logger

// MockTransactionLogger is a mock implementation of the TransactionLogger interface
type MockTransactionLogger struct {
	DeleteCalled    bool
	PutCalled       bool
	Events          chan Event
	ErrChannel      chan error
	Closed          bool
	WriteDeleteFunc func(key string) error
	WritePutFunc    func(key, value string) error
}

func (m *MockTransactionLogger) WriteDelete(key string) error {
	if m.WriteDeleteFunc == nil {
		return nil
	}
	return m.WriteDeleteFunc(key)
}

func (m *MockTransactionLogger) WritePut(key, value string) error {
	if m.WritePutFunc == nil {
		return nil
	}
	return m.WritePutFunc(key, value)
}

func (m *MockTransactionLogger) Err() <-chan error {
	return m.ErrChannel
}

func (m *MockTransactionLogger) Run() {

}

func (m *MockTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	return m.Events, m.ErrChannel
}

func (m *MockTransactionLogger) Close() error {
	m.Closed = true
	return nil
}
