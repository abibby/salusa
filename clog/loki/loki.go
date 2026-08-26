package loki

import (
	"log/slog"
	"net/url"

	"abibby.com/salusa/clog"
	"github.com/bearsoft-fi/slogloki"
)

type Config struct {
	Level    slog.Level
	URL      string
	TenantID string
}

var _ clog.Config = (*Config)(nil)

func (c *Config) Handler() (slog.Handler, error) {
	if _, err := url.ParseRequestURI(c.URL); err != nil {
		return nil, err
	}

	config := slogloki.NewDefaultConfig(
		c.URL,
		slogloki.WithLevel(c.Level),
		slogloki.WithTenantID(c.TenantID),
	)

	return slogloki.NewLokiHandler(config, map[string]string{})
}
