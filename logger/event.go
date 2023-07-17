package logger

type EventType byte

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)

type Event struct {
	Sequence int64
	Type     EventType
	Key      string
	Value    string
}
