package logger

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type FileTransactionLogger struct {
	events       chan Event
	errors       <-chan error
	lastSequence int64
	file         *os.File
}

func NewFileTransactionLogger(filename string) (*FileTransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}
	return &FileTransactionLogger{file: file}, nil
}

func (f *FileTransactionLogger) WritePut(key, value string) error {
	select {
	case f.events <- Event{
		Type:  EventPut,
		Key:   key,
		Value: value,
	}:
		return nil
	default:
		return errors.New("channel closed or nil")
	}
}

func (f *FileTransactionLogger) WriteDelete(key string) error {
	select {
	case f.events <- Event{
		Type: EventDelete,
		Key:  key,
	}:
		return nil
	default:
		return errors.New("channel closed or nil")
	}
}

func (f *FileTransactionLogger) Err() <-chan error {
	return f.errors
}

func (f *FileTransactionLogger) Run() {
	events := make(chan Event, 16)
	f.events = events

	errors := make(chan error)
	f.errors = errors

	go func() {
		for e := range events {
			f.lastSequence += 1
			_, err := fmt.Fprintf(f.file, "%d\t%d\t%s\t%s\n", f.lastSequence, e.Type, e.Key, e.Value)
			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

func (f *FileTransactionLogger) Close() error {
	if f.events != nil {
		close(f.events)
	}

	if err := f.file.Close(); err != nil {
		return err
	}

	return nil
}

func (f *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(f.file)
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		var e Event
		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			_, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s\n", &e.Sequence, &e.Type, &e.Key, &e.Value)
			if err != nil {
				_, err := fmt.Sscanf(line, "%d\t%d\t%s\n", &e.Sequence, &e.Type, &e.Key)
				e.Value = ""
				if err != nil {
					outError <- fmt.Errorf("input parse error: %w", err)
					return
				}
			}

			if f.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction number out of sequence")
				return
			}

			f.lastSequence = e.Sequence
			outEvent <- e

			if err := scanner.Err(); err != nil {
				outError <- fmt.Errorf("transaction log read failire: %w", err)
			}
		}
	}()
	return outEvent, outError
}
