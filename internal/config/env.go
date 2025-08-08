package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Env *schema

type schema struct {
	PORT              string
	RuntimeEnv        string
	POSTGRES_USER     string
	POSTGRES_DB       string
	POSTGRES_PASSWORD string
}

func Run(EnvFile string) *appConfig {
	cfg := &appConfig{
		EnvFile: EnvFile,
	}
	cfg.load()

	return cfg
}

type appConfig struct {
	EnvFile string
}

func (c *appConfig) load() {
	if err := godotenv.Load(c.EnvFile); err != nil {
		log.Fatalf("env file with %s error on loaded file\n err : %s", c.EnvFile, err)
	}

	Env = &schema{
		PORT:              *c.lookup("APP_PORTS"),
		RuntimeEnv:        *c.lookup("ENV"),
		POSTGRES_USER:     *c.lookup("POSTGRES_USER"),
		POSTGRES_DB:       *c.lookup("POSTGRES_DB"),
		POSTGRES_PASSWORD: *c.lookup("POSTGRES_PASSWORD"),
	}
}

func (c *appConfig) lookup(key string) *string {
	if value, ok := os.LookupEnv(key); ok {
		return &value
	}

	log.Fatalf("env value not available with key %s", key)
	return nil
}
