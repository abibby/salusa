package salusaconfig

import "context"

type Config interface {
	GetHTTPPort() int
	GetBaseURL() string
}

type ServiceConfig interface {
	Register(ctx context.Context) error
}
