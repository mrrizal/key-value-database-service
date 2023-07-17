package logger

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error
	Run()
	ReadEvents() (<-chan Event, <-chan error)
	Close() error
}
