package repository

import (
	"context"
	"fmt"
	"strukit-services/pkg/logger"
	"time"

	"gorm.io/gorm"
)

func NewBase(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

type BaseRepository struct {
	db *gorm.DB
}

func (b *BaseRepository) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	tx := b.db.Begin()
	defer func() {
		if err := recover(); err != nil {
			logger.Log.DB(ctx).Errorf("error with transaction, err:%s", err)
			tx.Rollback()
			panic(err)
		}
	}()

	if err := fn(tx); err != nil {
		if re := tx.Rollback().Error; re != nil {
			return fmt.Errorf("transaction failed: %v, rollback failed: %w", err, re)
		}

		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed commit transaction %w", err)
	}

	return nil
}

func (b BaseRepository) Now() *time.Time {
	n := time.Now()
	return &n
}
