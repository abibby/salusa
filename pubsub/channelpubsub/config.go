package channelpubsub

import (
	"context"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/pubsub"
)

func Register(ctx context.Context) {
	di.RegisterLazySingleton(ctx, func() (pubsub.PubSub, error) {
		return New(), nil
	})
	pubsub.RegisterTopic(ctx)
}
