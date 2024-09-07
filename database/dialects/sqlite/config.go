package sqlite

import (
	"github.com/abibby/salusa/database"
	// _ "modernc.org/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Path string
}

var _ database.Config = (*Config)(nil)

func NewConfig(path string) *Config {
	return &Config{
		Path: path,
	}
}

func (c *Config) SetDialect() {
	UseSQLite()
}
func (c *Config) DriverName() string {
	return "sqlite3"
}
func (c *Config) DataSourceName() string {
	return c.Path
}
