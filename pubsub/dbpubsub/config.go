package dbpubsub

import (
	"context"

	"abibby.com/salusa/database"
	"abibby.com/salusa/di"
	"abibby.com/salusa/pubsub"
)

func Register(ctx context.Context) {
	di.RegisterLazySingletonWith(ctx, func(u database.Update) (pubsub.PubSub, error) {
		return New(u), nil
	})
	pubsub.RegisterTopic(ctx)
}
