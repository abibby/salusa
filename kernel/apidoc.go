package kernel

import (
	"context"
	"net/url"

	"github.com/abibby/salusa/openapidoc"
	"github.com/go-openapi/spec"
)

// var _ openapidoc.APIDocer = (*Kernel)(nil)

// Operation implements openapidoc.Operationer.
func (k *Kernel) APIDoc(ctx context.Context) (*spec.Swagger, error) {
	var docs spec.Swagger

	if k.docs != nil {
		docs = *k.docs
	}
	h := k.RootHandler(context.Background())
	if paths, ok := h.(openapidoc.Pathser); ok {
		var err error
		docs.Paths, err = paths.Paths()
		if err != nil {
			return nil, err
		}
	}

	u, err := url.Parse(k.cfg.GetBaseURL())
	if err != nil {
		k.Logger(context.Background()).Warn("invalid base url")
	} else {
		docs.Schemes = []string{u.Scheme}
		docs.Host = u.Host
		docs.BasePath = u.Path
	}

	docs.Swagger = "2.0"

	return &docs, nil
}
