package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvWithDefault(t *testing.T) {
	t.Run("returns value from environment variable", func(t *testing.T) {
		// Устанавливаем переменную окружения
		os.Setenv("TEST_VAR", "test_value")
		defer os.Unsetenv("TEST_VAR") // Очищаем переменную окружения после теста

		// Проверяем, что вернется значение из окружения
		result := GetEnvWithDefault("TEST_VAR", "default_value")
		assert.Equal(t, "test_value", result)
	})

	t.Run("returns default value if environment variable is not set", func(t *testing.T) {
		// Проверяем, что вернется значение по умолчанию
		result := GetEnvWithDefault("NON_EXISTENT_VAR", "default_value")
		assert.Equal(t, "default_value", result)
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("loads default config when environment variables are not set", func(t *testing.T) {
		// Очищаем все переменные окружения
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")

		// Загружаем конфиг
		config := LoadConfig()

		// Проверяем, что значения конфигурации по умолчанию установлены правильно
		assert.Equal(t, "8080", config.Port)
		assert.Equal(t, "ttavito", config.DBUser)
		assert.Equal(t, "ttavito", config.DBPassword)
		assert.Equal(t, "localhost", config.DBHost)
		assert.Equal(t, "5432", config.DBPort)
		assert.Equal(t, "ttavito", config.DBName)
	})

	t.Run("loads config from environment variables", func(t *testing.T) {
		// Устанавливаем переменные окружения для теста
		os.Setenv("DB_USER", "custom_user")
		os.Setenv("DB_PASSWORD", "custom_password")
		os.Setenv("DB_HOST", "custom_host")
		os.Setenv("DB_PORT", "1234")
		os.Setenv("DB_NAME", "custom_db")
		defer func() {
			// Очищаем переменные окружения после теста
			os.Unsetenv("DB_USER")
			os.Unsetenv("DB_PASSWORD")
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_PORT")
			os.Unsetenv("DB_NAME")
		}()

		// Загружаем конфиг
		config := LoadConfig()

		// Проверяем, что значения конфигурации загружены из окружения
		assert.Equal(t, "8080", config.Port)
		assert.Equal(t, "custom_user", config.DBUser)
		assert.Equal(t, "custom_password", config.DBPassword)
		assert.Equal(t, "custom_host", config.DBHost)
		assert.Equal(t, "1234", config.DBPort)
		assert.Equal(t, "custom_db", config.DBName)
	})
}
