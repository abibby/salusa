package postgres

import (
	"fmt"

	"github.com/abibby/salusa/database"
	_ "github.com/lib/pq"
)

type Config struct {
	Username   string
	Password   string
	Host       string
	Database   string
	DisableSSL bool
}

var _ database.Config = (*Config)(nil)

func (c *Config) SetDialect() {
	UsePostgres()
}
func (c *Config) DriverName() string {
	return "postgres"
}
func (c *Config) DataSourceName() string {
	ssl := ""
	if c.DisableSSL {
		ssl = "sslmode=disable"
	}
	return fmt.Sprintf("host=%s dbname=%s user=%s password=%s %s", c.Host, c.Database, c.Username, c.Password, ssl)
}
