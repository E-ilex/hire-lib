package config

import "os"

const (
	envPortName          = "PORT"
	envPortValue         = "3000"
	envDBDriverName      = "DB_DRIVER_HTL"
	envDBDriverValue     = "sqlite3"
	envDBDatasourceName  = "DB_DATASOURCE_HTL"
	envDBDatasourceValue = "qa.sqlite3"
)

type Configuration struct {
	Port         string
	DBDriver     string
	DBDatasource string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func LoadConfiguration() *Configuration {
	port := getEnv(envPortName, envPortValue)
	dbDriver := getEnv(envDBDriverName, envDBDriverValue)
	dbDatasource := getEnv(envDBDatasourceName, envDBDatasourceValue)

	return &Configuration{
		Port:         port,
		DBDriver:     dbDriver,
		DBDatasource: dbDatasource,
	}
}
