package channelpubsub_test

import (
	"testing"

	"github.com/abibby/salusa/pubsub"
	"github.com/abibby/salusa/pubsub/channelpubsub"
	"github.com/abibby/salusa/pubsub/pubsubtest"
)

func TestStandard(t *testing.T) {
	pubsubtest.RunStandardTests(t, func(t *testing.T, name string, fn func(*testing.T, pubsub.PubSub)) {
		t.Run(name, func(t *testing.T) {
			p := channelpubsub.New()
			fn(t, p)
		})
	})
}
