package dbpubsub

import (
	"context"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/pubsub"
)

func Register(ctx context.Context, table string) {
	di.RegisterLazySingletonWith(ctx, func(u database.Update) (pubsub.PubSub, error) {
		return New(u), nil
	})
	pubsub.RegisterTopic(ctx)
}
