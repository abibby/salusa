package database

type DBConfiger interface {
	DBConfig() Config
}

type Config interface {
	DriverName() string
	DataSourceName() string
	SetDialect()
}
