package db

import (
	"fmt"
	"strukit-services/internal/config"
)

func DsnPg() string {
	cfg := &dbCofig{}
	return cfg.BuildPgDsn()
}

type dbCofig struct {
	Host     string
	Port     string
	DbName   string
	Username string
	Password string
	SSLMode  bool
	TimeZone *string
}

func (c *dbCofig) SetSSLMode(value bool) {
	c.SSLMode = value
}

func (c *dbCofig) buildConfig() {
	timeZone := "Asia/Jakarta"
	c.Host = config.Env.DB_HOST
	c.Port = config.Env.DB_PORT
	c.DbName = config.Env.POSTGRES_DB
	c.Username = config.Env.POSTGRES_USER
	c.Password = config.Env.POSTGRES_PASSWORD
	c.SSLMode = false
	c.TimeZone = &timeZone
}

func (c *dbCofig) BuildPgDsn() string {
	c.buildConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=%s", c.Host, c.Port, c.Username, c.Password, c.DbName, *c.TimeZone)
	return dsn
}
