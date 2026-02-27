package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type EnvProvider struct {
	dbConnectionString     string
	valkeyConnectionString string
}

func NewEnvProvider() *EnvProvider {

	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	valkeyConnectionString := os.Getenv("VALKEY_CONNECTION_STRING")

	return &EnvProvider{
		dbConnectionString:     dbConnectionString,
		valkeyConnectionString: valkeyConnectionString,
	}
}

func (p *EnvProvider) GetDBConnectionString() string {
	return p.dbConnectionString
}

func (p *EnvProvider) GetValkeyConnectionString() string {
	return p.valkeyConnectionString
}
