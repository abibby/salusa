package loki

import (
	"log/slog"

	"github.com/abibby/salusa/clog"
	"github.com/grafana/loki-client-go/loki"
	slogloki "github.com/samber/slog-loki/v3"
)

type Config struct {
	URL      string
	TenantID string
}

var _ clog.Config = (*Config)(nil)

func (c *Config) Handler() (slog.Handler, error) {

	// setup loki client
	config, err := loki.NewDefaultConfig(c.URL)
	if err != nil {
		return nil, err
	}

	config.TenantID = c.TenantID

	// client, err := loki.NewWithLogger(config, &localLogger{slog.New(clog.DefaultHandler())})
	client, err := loki.New(config)
	if err != nil {
		return nil, err
	}

	return slogloki.Option{Level: slog.LevelDebug, Client: client}.NewLokiHandler(), nil
}
