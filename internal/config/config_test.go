package config

import (
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testDataFolder = "testdata/"

func init() {
	log.SetOutput(io.Discard)
}

func TestLoadConfig(t *testing.T) {
	t.Run("non_existent_file", func(t *testing.T) {
		_, err := LoadConfig("dfklajflkjd")
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("load_config_yaml", func(t *testing.T) {
		expected := &Config{
			Server: ServerConfig{
				Host:            "127.0.0.1",
				Port:            5000,
				ShutdownTimeout: 10 * time.Second,
				IsProduction:    true,
				Addr:            "127.0.0.1:5000",
			},
			Log: LogConfig{Level: "debug"},
			Security: SecurityConfig{
				PasswordSalt:     "something",
				JWTSecret:        "amazing",
				AccessExpiresIn:  time.Minute * 5,
				RefreshExpiresIn: time.Hour * 24 * 30,
			},
			Postgres: PostgresConfig{
				Username: "postgre",
				Password: "postgre",
				Host:     "127.0.0.1",
				Port:     2345,
				Database: "postgre",
				URI:      "postgres://postgre:postgre@127.0.0.1:2345/postgre",
			},
		}
		actual, err := LoadConfig(testDataFolder + "correct.yaml")

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
