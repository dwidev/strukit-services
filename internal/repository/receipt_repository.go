package repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strukit-services/internal/models"
	appContext "strukit-services/pkg/context"
	"strukit-services/pkg/responses"

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

func (r *ReceiptRepository) GetDetailReceipt(ctx context.Context) (receipts *models.Receipt, err error) {
	receiptId := ctx.Value(appContext.ReceiptIDKey).(uuid.UUID)
	query := r.db.WithContext(ctx).Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, receipt_id, item_name, quantity, unit_price, total_price, discount_amount")
	})

	if err = query.Where("id = ?", receiptId).First(&receipts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, responses.Err(http.StatusNotFound, fmt.Sprintf("not found receipt with id %s", receiptId))
		}
		return nil, err
	}

	return
}

func (r *ReceiptRepository) GetReceiptByProjectID(ctx context.Context) (receipts []*models.Receipt, err error) {
	projectId := ctx.Value(appContext.ProjectID).(uuid.UUID)
	query := r.db.WithContext(ctx).Model(&models.Receipt{}).Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, receipt_id, item_name, quantity, unit_price, total_price, discount_amount")
	})

	if err = query.Where("project_id = ?", projectId).Find(&receipts).Error; err != nil {
		return nil, err
	}

	return
}

func (r *ReceiptRepository) Delete(ctx context.Context) error {
	receiptID := ctx.Value(appContext.ReceiptIDKey).(uuid.UUID)
	if err := r.db.Model(&models.Receipt{}).Delete("id = ?", receiptID).Error; err != nil {
		return err
	}
	return nil
}

func (r *ReceiptRepository) FindByFingerprint(ctx context.Context, fingerprint string) ([]*models.Receipt, error) {
	projectId := ctx.Value(appContext.ProjectID).(uuid.UUID)
	var receipts []*models.Receipt
	if err := r.db.WithContext(ctx).Where("fingerprint = ? AND project_id = ?", fingerprint, projectId).Find(&receipts).Error; err != nil {
		return nil, err
	}

	return receipts, nil
}

func (r *ReceiptRepository) FindByContentHash(ctx context.Context, hash string) ([]*models.Receipt, error) {
	projectId := ctx.Value(appContext.ProjectID).(uuid.UUID)
	var receipts []*models.Receipt
	if err := r.db.WithContext(ctx).Where("content_hash = ? AND project_id = ?", hash, projectId).Find(&receipts).Error; err != nil {
		return nil, err
	}

	return receipts, nil
}

func (r *ReceiptRepository) FindSimilarReceipts(ctx context.Context, criteria *models.Receipt) ([]*models.Receipt, error) {
	var receipts []*models.Receipt

	query := r.db.WithContext(ctx).Where("project_id = ?", criteria.ProjectID)

	if criteria.MerchantName != nil {
		query = query.Where("merchant_name ILIKE ?", "%"+*criteria.MerchantName+"%")
	}

	minAmount := criteria.TotalAmount * (1 - 0.2)
	maxAmount := criteria.TotalAmount * (1 + 0.3)
	query = query.Where("total_amount BETWEEN ? AND ?", minAmount, maxAmount)

	startDate := criteria.TransactionDate.AddDate(0, 0, -1)
	endDate := criteria.TransactionDate.AddDate(0, 0, 1)
	query = query.Where("transaction_date BETWEEN ? AND ?", startDate, endDate)

	if err := query.Find(&receipts).Error; err != nil {
		return nil, err
	}

	return receipts, nil
}
