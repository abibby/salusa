package filesystem

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/salusaconfig"
)

type Config interface {
	FS() fs.FS
}

type FSConfiger interface {
	FSConfig() Config
}

type LocalFS struct {
	Root string
}

func NewLocalFS(root string) *LocalFS {
	return &LocalFS{
		Root: root,
	}
}

func (l *LocalFS) FS() fs.FS {
	return os.DirFS(l.Root)
}

func Register(ctx context.Context) error {
	di.RegisterLazySingletonWith(ctx, func(cfg salusaconfig.Config) (fs.FS, error) {
		var cfgAny any = cfg
		cfger, ok := cfgAny.(FSConfiger)
		if !ok {
			return nil, fmt.Errorf("config not instance of email.FSConfiger")
		}
		return cfger.FSConfig().FS(), nil
	})
	return nil
}
