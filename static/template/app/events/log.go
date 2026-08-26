package events

import (
	"abibby.com/salusa/event"
	"abibby.com/salusa/event/cron"
)

type LogEvent struct {
	cron.CronEvent
	Message string
}

var _ event.Event = (*LogEvent)(nil)

func (e *LogEvent) Type() event.EventType {
	return "template:example-event"
}
