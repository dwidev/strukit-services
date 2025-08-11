package repository

import (
	"gorm.io/gorm"
)

func NewBase(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

type BaseRepository struct {
	db *gorm.DB
}
