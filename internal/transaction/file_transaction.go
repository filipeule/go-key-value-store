package transaction

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type FileTransactionLogger struct {
	events       chan<- Event
	errors       <-chan error
	lastSequence uint64
	file         *os.File
}

func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}

	return &FileTransactionLogger{file: file}, nil
}

func (pl *FileTransactionLogger) Run() {
	events := make(chan Event, 16)
	pl.events = events

	errors := make(chan error, 1)
	pl.errors = errors

	go func() {
		for e := range events {
			pl.lastSequence++
			_, err := fmt.Fprintf(
				pl.file,
				"%d\t%d\t%s\t%s\n",
				pl.lastSequence, e.EventType, e.Key, e.Value,
			)

			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

func (pl *FileTransactionLogger) WritePut(key, value string) {
	pl.events <- Event{EventType: EventPut, Key: key, Value: value}
}

func (pl *FileTransactionLogger) WriteDelete(key string) {
	pl.events <- Event{EventType: EventDelete, Key: key}
}

func (pl *FileTransactionLogger) Err() <-chan error {
	return pl.errors
}

func (pl *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(pl.file)

	outEvent := make(chan Event)
	outError := make(chan error)

	go func() {
		var e Event

		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			parts := strings.SplitN(line, "\t", 4)
			if len(parts) < 3 {
				outError <- fmt.Errorf("invalid logger line (too few fields): %q", line)
				return
			}

			seq, err := strconv.ParseUint(parts[0], 10, 64)
			if err != nil {
				outError <- fmt.Errorf("invalid sequence: %w", err)
				return
			}

			eventTypeInt, err := strconv.Atoi(parts[1])
			if err != nil {
				outError <- fmt.Errorf("invalid event type: %w", err)
				return
			}

			e = Event{
				Sequence:  seq,
				EventType: EventType(eventTypeInt),
				Key:       parts[2],
			}

			if e.EventType == EventPut {
				if len(parts) != 4 {
					outError <- fmt.Errorf("PUT event missing value field: %q", line)
					return
				}
				e.Value = parts[3]
			}

			if pl.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			pl.lastSequence = e.Sequence
			outEvent <- e
		}

		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()

	return outEvent, outError
}
