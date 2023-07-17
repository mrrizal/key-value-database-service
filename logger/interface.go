package logger

type TransactionLogger interface {
	WriteDelete(key string) error
	WritePut(key, value string) error
	Err() <-chan error
	Run()
	ReadEvents() (<-chan Event, <-chan error)
	Close() error
}
