package repository

import (
	"context"
	"fmt"
	appContext "strukit-services/pkg/context"
	"strukit-services/pkg/logger"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func NewBase(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

type BaseRepository struct {
	db *gorm.DB
}

func (b *BaseRepository) UserID(ctx context.Context) uuid.UUID {
	userId := ctx.Value(appContext.UserIDKey).(uuid.UUID)
	return userId
}

func (b *BaseRepository) ProjectID(ctx context.Context) uuid.UUID {
	projectId := ctx.Value(appContext.ProjectID).(uuid.UUID)
	return projectId
}

func (b *BaseRepository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
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

// function for get date and time at now
func (b BaseRepository) Now() *time.Time {
	n := time.Now()
	return &n
}

// function for get date now without the time
func (b BaseRepository) DateNow() *time.Time {
	n := time.Now()
	date := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location())
	return &date
}
