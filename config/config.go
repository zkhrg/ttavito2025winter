package config

type Config struct {
	DatabaseURL string
	Port        string
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL: "postgresql://ttavito:ttavito@postgres/ttavito?sslmode=disable",
		Port:        "8080",
	}
}
