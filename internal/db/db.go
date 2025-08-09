package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open() *Database {
	dsn := DsnPg()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Fatalf("error with open database, err : %s", err)
	}

	return &Database{
		db: db,
	}
}

type Database struct {
	db *gorm.DB
}

func (d *Database) Instance() *gorm.DB {
	return d.db
}

func (d *Database) Close() {
	sqlDb, err := d.db.DB()
	if err != nil {
		log.Fatalf("error open sqldb, err : %s", err)
	}
	sqlDb.Close()
}
