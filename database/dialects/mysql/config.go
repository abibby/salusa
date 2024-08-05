package mysql

import (
	"github.com/abibby/salusa/database"
	"github.com/go-sql-driver/mysql"
)

type SimpleConfig struct {
	Username string
	Password string
	Address  string
	Database string
}

var _ database.Config = (*SimpleConfig)(nil)

func (c *SimpleConfig) SetDialect() {
	UseMySql()
}
func (c *SimpleConfig) DriverName() string {
	return "mysql"
}
func (c *SimpleConfig) DataSourceName() string {
	mysqlCfg := mysql.NewConfig()
	mysqlCfg.User = c.Username
	mysqlCfg.Passwd = c.Password
	mysqlCfg.Addr = c.Address
	mysqlCfg.DBName = c.Database
	return mysqlCfg.FormatDSN()
}

type Config struct {
	cfg *mysql.Config
}

func NewMySQLConfig(cfg *mysql.Config) *Config {
	return &Config{
		cfg: cfg,
	}
}

var _ database.Config = (*Config)(nil)

func (c *Config) SetDialect() {
	UseMySql()
}
func (c *Config) DriverName() string {
	return "mysql"
}
func (c *Config) DataSourceName() string {
	return c.cfg.FormatDSN()
}
