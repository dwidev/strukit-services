package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Env *schema

var (
	Prod = "production"
	Dev  = "dev"
)

type schema struct {
	PORT              string
	RuntimeEnv        string
	POSTGRES_USER     string
	POSTGRES_DB       string
	POSTGRES_PASSWORD string
}

func Run() *appConfig {
	env := os.Getenv("GO_ENV")
	cfg := &appConfig{
		Env: env,
	}
	cfg.load()

	return cfg
}

type appConfig struct {
	Env string
}

func (c *appConfig) load() {
	var file string

	switch c.Env {
	case "dev":
		file = ".env.dev"
	default:
		file = ".env"
	}

	log.Printf("load env file with %s env\n", c.Env)
	if err := godotenv.Load(file); err != nil {
		log.Fatalf("env file with %s error on loaded file\n err : %s", file, err)
	}

	Env = &schema{
		PORT:              *c.lookup("APP_PORTS"),
		RuntimeEnv:        c.Env,
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
