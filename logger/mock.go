package logger

import (
	"fmt"
)

// MockTransactionLogger is a mock implementation of the TransactionLogger interface
type MockTransactionLogger struct {
	DeleteCalled bool
	PutCalled    bool
	Events       chan Event
	ErrChannel   chan error
	Closed       bool
}

func (m *MockTransactionLogger) WriteDelete(key string) {
	m.DeleteCalled = true
	fmt.Println("WriteDelete called with key:", key)
}

func (m *MockTransactionLogger) WritePut(key, value string) {
	m.PutCalled = true
	fmt.Println("WritePut called with key:", key, "and value:", value)
}

func (m *MockTransactionLogger) Err() <-chan error {
	return m.ErrChannel
}

func (m *MockTransactionLogger) Run() {
	fmt.Println("Run called")
}

func (m *MockTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	return m.Events, m.ErrChannel
}

func (m *MockTransactionLogger) Close() error {
	m.Closed = true
	fmt.Println("Close called")
	return nil
}
