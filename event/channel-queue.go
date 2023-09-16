package event

type ChannelQueue struct {
	channel chan Event
}

func NewChannelQueue() *ChannelQueue {
	return &ChannelQueue{
		channel: make(chan Event, 10),
	}
}

func (q ChannelQueue) Push(e Event) {
	q.channel <- e
}
func (q ChannelQueue) Pop() Event {
	return <-q.channel
}
