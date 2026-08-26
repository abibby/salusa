package channelpubsub

import (
	"context"

	"abibby.com/salusa/di"
	"abibby.com/salusa/pubsub"
)

func Register(ctx context.Context) {
	di.RegisterLazySingleton(ctx, func() (pubsub.PubSub, error) {
		return New(), nil
	})
	pubsub.RegisterTopic(ctx)
}
