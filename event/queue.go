package event

type Queue interface {
	Push(e Event)
	Pop() Event
}
