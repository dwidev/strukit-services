package repository

import (
	"context"
	"strukit-services/internal/models"
	appContext "strukit-services/pkg/context"

	"github.com/google/uuid"
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
	data.UserID = ctx.Value(appContext.UserIDKey).(uuid.UUID)
	data.ProjectID = ctx.Value(appContext.ProjectID).(uuid.UUID)
	if err := r.db.WithContext(ctx).Save(data).Error; err != nil {
		return nil, err
	}

	return data, nil
}
