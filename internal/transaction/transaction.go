package transaction

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)

const (
	PostgresTransaction TransactionType = "postgres"
	FileTransaction TransactionType = "file"
)

type TransactionType string

type EventType byte

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error
	ReadEvents() (<-chan Event, <-chan error)
	Run()
}

type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     string
}
