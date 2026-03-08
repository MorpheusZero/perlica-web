package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type EnvProvider struct {
	dbConnectionString     string
	valkeyConnectionString string
	hostDomain             string
}

func NewEnvProvider() *EnvProvider {

	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	valkeyConnectionString := os.Getenv("VALKEY_CONNECTION_STRING")
	hostDomain := os.Getenv("HOST_DOMAIN")

	return &EnvProvider{
		dbConnectionString:     dbConnectionString,
		valkeyConnectionString: valkeyConnectionString,
		hostDomain:             hostDomain,
	}
}

func (p *EnvProvider) GetDBConnectionString() string {
	return p.dbConnectionString
}

func (p *EnvProvider) GetValkeyConnectionString() string {
	return p.valkeyConnectionString
}

func (p *EnvProvider) GetHostDomain() string {
	return p.hostDomain
}
