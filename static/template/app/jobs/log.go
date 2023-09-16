package jobs

import (
	"log"

	"github.com/abibby/salusa/static/template/app/events"
)

func LogJob(e *events.LogEvent) error {
	log.Print(e.Message)
	return nil
}
