package salusaconfig

type Config interface {
	GetHTTPPort() int
	GetBaseURL() string
}
