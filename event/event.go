package event

type EventType string

type Event interface {
	Type() EventType
}
