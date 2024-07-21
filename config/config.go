package config

type Config interface {
	GetHTTPPort() int
	GetBaseURL() string
}
