package repository

import (
	"context"
	"strukit-services/internal/models"
	appContext "strukit-services/pkg/context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func NewReceipt(base *BaseRepository) *ReceiptRepository {
	return &ReceiptRepository{
		BaseRepository: base,
	}
}

type ReceiptRepository struct {
	*BaseRepository
}

func (r *ReceiptRepository) Save(ctx context.Context, data *models.Receipt) (*models.Receipt, error) {
	err := r.Transaction(ctx, func(tx *gorm.DB) error {
		data.UserID = ctx.Value(appContext.UserIDKey).(uuid.UUID)
		data.ProjectID = ctx.Value(appContext.ProjectID).(uuid.UUID)

		items := data.Items
		data.Items = nil
		if err := tx.WithContext(ctx).Save(data).Error; err != nil {
			return err
		}

		if len(items) > 0 {
			if err := tx.WithContext(ctx).CreateInBatches(items, len(items)).Error; err != nil {
				return err
			}
		}

		data.Items = items
		return nil
	})

	if err != nil {
		return nil, err
	}

	return data, nil
}
