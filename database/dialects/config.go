package dialects

type Config interface {
	DriverName() string
	DataSourceName() string
	SetDialect()
}
