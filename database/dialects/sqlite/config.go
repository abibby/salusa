package sqlite

import (
	// _ "modernc.org/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Path string
}

func NewConfig(path string) *Config {
	return &Config{
		Path: path,
	}
}

func (c *Config) DriverName() string {
	return "sqlite3"
}
func (c *Config) DataSourceName() string {
	return c.Path
}
