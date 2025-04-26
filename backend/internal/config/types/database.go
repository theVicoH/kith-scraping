package types

type DatabaseConfig struct {
	Name     string `env:"NAME"`
	Host     string `env:"HOST"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Port     int    `env:"PORT"`
}