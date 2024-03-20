package sqlite

import (
	_ "modernc.org/sqlite"
)

type Config struct {
	Path string
}

func NewConfig(path string) *Config {
	return &Config{
		Path: path,
	}
}

func (c *Config) SetDialect() {
	UseSQLite()
}
func (c *Config) DriverName() string {
	return "sqlite"
}
func (c *Config) DataSourceName() string {
	return c.Path
}
