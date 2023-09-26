package event

import "reflect"

type ChannelQueue struct {
	channel chan []byte
}

func NewChannelQueue() *ChannelQueue {
	return &ChannelQueue{
		channel: make(chan []byte, 10),
	}
}

func (q ChannelQueue) Push(e Event) error {
	b, err := encodeEvent(e)
	if err != nil {
		return err
	}
	q.channel <- b
	return nil
}
func (q ChannelQueue) Pop(events map[EventType]reflect.Type) (Event, error) {
	return decodeEvent(<-q.channel, events)
}
