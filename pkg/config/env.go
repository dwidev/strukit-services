package config

import (
	"os"
	"strukit-services/pkg/constant"
	"strukit-services/pkg/logger"

	"github.com/joho/godotenv"
)

var Env *schema

type schema struct {
	RuntimeEnv constant.Environment
	PORT       string

	DB_HOST string
	DB_PORT string

	POSTGRES_DB       string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string

	JWT_ACCESS_SECRET  string
	JWT_REFRESH_SECRET string
}

func Run(env string) *appConfig {
	if env == "" {
		env = string(constant.Dev)
	}

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
	file := c.GetEnvFile()
	if err := godotenv.Load(file); err != nil {
		logger.Log.Fatalf("env file with %s error on loaded file\n err : %s", file, err)

	}

	schema := &schema{
		PORT:       *c.lookup("APP_PORTS"),
		RuntimeEnv: constant.Environment(c.Env),

		DB_HOST: *c.lookup("DB_HOST"),
		DB_PORT: *c.lookup("DB_PORT"),

		POSTGRES_DB:       *c.lookup("POSTGRES_DB"),
		POSTGRES_USER:     *c.lookup("POSTGRES_USER"),
		POSTGRES_PASSWORD: *c.lookup("POSTGRES_PASSWORD"),

		JWT_ACCESS_SECRET:  *c.lookup("JWT_ACCESS_SECRET"),
		JWT_REFRESH_SECRET: *c.lookup("JWT_REFRESH_SECRET"),
	}

	Env = schema
}

func (c *appConfig) lookup(key string) *string {
	if value, ok := os.LookupEnv(key); ok {
		return &value
	}

	logger.Log.Fatalf("env value not available with key %s", key)
	return nil
}

func (c *appConfig) GetEnvFile() string {
	var file string

	switch c.Env {
	case "dev":
		file = ".env.dev"
	default:
		file = ".env"
	}

	return file
}
