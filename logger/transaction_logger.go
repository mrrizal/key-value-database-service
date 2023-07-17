package logger

import (
	"fmt"
	"log"
)

func InitializeTransactionLog(name string) (TransactionLogger, error) {
	var err error

	logger, err := NewFileTransactionLogger(name)
	if err != nil {
		return nil, fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := logger.ReadEvents()
	e, ok := Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.Type {
			case EventDelete:
				log.Println("delete event")
			case EventPut:
				log.Println("put event")
			}
		}
	}

	logger.Run()
	return logger, nil
}
