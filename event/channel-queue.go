package event

import "github.com/barkimedes/go-deepcopy"

type ChannelQueue struct {
	channel chan Event
}

func NewChannelQueue() *ChannelQueue {
	return &ChannelQueue{
		channel: make(chan Event, 10),
	}
}

func (q ChannelQueue) Push(e Event) {
	q.channel <- deepcopy.MustAnything(e).(Event)
}
func (q ChannelQueue) Pop() Event {
	return <-q.channel
}
