package pubsub_test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/pubsub"
)

// mockTopic implements pubsub.Topic.
type mockTopic struct {
	name string
}

func (t *mockTopic) Publish(ctx context.Context, data []byte) {}
func (t *mockTopic) Subscribe(ctx context.Context) pubsub.Subscription {
	return &mockSubscription{}
}

// mockSubscription implements pubsub.Subscription.
type mockSubscription struct{}

func (s *mockSubscription) Close() error { return nil }
func (s *mockSubscription) Next() []byte { return []byte("hello") }

// mockPubSub implements pubsub.PubSub.
type mockPubSub struct{}

func (p *mockPubSub) Topic(name string) pubsub.Topic {
	return &mockTopic{name: name}
}

// Compile-time assertions that the mock types satisfy the interfaces.
var (
	_ pubsub.PubSub       = (*mockPubSub)(nil)
	_ pubsub.Topic        = (*mockTopic)(nil)
	_ pubsub.Subscription = (*mockSubscription)(nil)
)

func TestPubSubInterfaces(t *testing.T) {
	ps := &mockPubSub{}
	sub := ps.Topic("test").Subscribe(context.Background())
	if sub.Next() == nil {
		t.Fatal("expected non-nil data")
	}
	if err := sub.Close(); err != nil {
		t.Fatal(err)
	}
}
