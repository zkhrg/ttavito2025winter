package config

type Config struct {
	DatabaseURL string
	Port        string
}

func LoadConfig() *Config {
	return &Config{
		// DatabaseURL: os.Getenv("DATABASE_URL"),
		// Port:        os.Getenv("PORT"),

		DatabaseURL: "postgresql://gitverse-internship-zg:gitverse-internship-zg@postgres/gitverse-internship-zg?sslmode=disable",
		Port:        "8080",
	}
}
