package mysql

import (
	"github.com/go-sql-driver/mysql"
)

type SimpleConfig struct {
	Username string
	Password string
	Host     string
	Database string
}

func (c *SimpleConfig) DriverName() string {
	return "mysql"
}
func (c *SimpleConfig) DataSourceName() string {
	mysqlCfg := mysql.NewConfig()
	mysqlCfg.User = c.Username
	mysqlCfg.Passwd = c.Password
	mysqlCfg.Addr = c.Host
	mysqlCfg.DBName = c.Database
	mysqlCfg.MultiStatements = true
	mysqlCfg.ParseTime = true
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
func (c *Config) DriverName() string {
	return "mysql"
}
func (c *Config) DataSourceName() string {
	return c.cfg.FormatDSN()
}

// func (c *Config) Register(ctx context.Context) error {
// 	di.RegisterLazySingleton(ctx, func() (*sqlx.DB, error) {
// 		UseMySql()
// 		return sqlx.Open(c.DriverName(), c.DataSourceName())
// 	})

// 	return databasedi.RegisterTransactions(nil)(ctx)
// }
