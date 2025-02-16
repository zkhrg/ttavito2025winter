package config

import "os"

type Config struct {
	Port       string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

func GetEnvWithDefault(key string, defaultValue string) string {
	// Получаем значение переменной окружения
	value := os.Getenv(key)

	// Если переменная окружения не установлена (пустое значение), возвращаем значение по умолчанию
	if value == "" {
		return defaultValue
	}
	return value
}

func LoadConfig() *Config {
	return &Config{
		Port:       "8080",
		DBUser:     GetEnvWithDefault("DB_USER", "ttavito"),
		DBPassword: GetEnvWithDefault("DB_PASSWORD", "ttavito"),
		DBHost:     GetEnvWithDefault("DB_HOST", "localhost"),
		DBPort:     GetEnvWithDefault("DB_PORT", "5432"),
		DBName:     GetEnvWithDefault("DB_NAME", "ttavito"),
	}
}
