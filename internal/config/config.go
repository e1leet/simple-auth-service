package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Host            string        `yaml:"host" env-default:"localhost"`
	Port            int           `yaml:"port" env-default:"8000"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"5s"`
	IsProduction    bool          `yaml:"is_production" env-default:"false"`
	Addr            string
}

type LogConfig struct {
	Level string `yaml:"level" env-default:"trace"`
}

type SecurityConfig struct {
	PasswordSalt     string        `yaml:"password_salt" env-default:"salt"`
	JWTSecret        string        `yaml:"jwt_secret" env-default:"secret"`
	AccessExpiresIn  time.Duration `yaml:"access_expires_in" env-default:"30m"`
	RefreshExpiresIn time.Duration `yaml:"refresh_expires_in" env-default:"1440h"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Log      LogConfig      `yaml:"log"`
	Security SecurityConfig `yaml:"security"`
}

func (c *Config) Update() error {
	c.Server.Addr = fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
	return nil
}

func LoadConfig(path string) (*Config, error) {
	log.Printf("config init. path=%s\n", path)

	instance := &Config{}

	if err := cleanenv.ReadConfig(path, instance); err != nil {
		log.Println(err)
		return nil, err
	}

	return instance, nil
}
