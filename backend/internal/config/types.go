package config

// Config represents the application configuration

// DatabaseConfig represents database-specific configuration
type DatabaseConfig struct {
	Name     string `env:"NAME" required:"true"`
	Host     string `env:"HOST required:"true"`
	User     string `env:"USER required:"true"`
	Password string `env:"PASSWORD required:"true"`
	Port     int    `env:"PORT required:"true"`
}
