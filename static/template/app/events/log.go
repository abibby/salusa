package events

import (
	"github.com/abibby/salusa/event"
	"github.com/abibby/salusa/event/cron"
)

type LogEvent struct {
	cron.BaseEvent
	Message string
}

var _ event.Event = (*LogEvent)(nil)

func (e *LogEvent) Type() event.EventType {
	return "template:example-event"
}

// func init() {
// 	kernel.RegisterEvent(&LogEvent{})
// }
