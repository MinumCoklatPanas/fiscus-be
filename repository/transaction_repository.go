package repository

import (
	"context"
	"fiscus/model"
	"fiscus/util"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TrxRepository interface {
	FindAll(
		ctx context.Context,
		trxType model.TrxTypeEnum,
		startDate time.Time,
		endDate time.Time,
		accountId uuid.UUID,
		categoryId uuid.UUID,
		pagination util.Pagination,
	) ([]model.Transaction, int64, error)

	Create(ctx context.Context, trx model.Transaction) (model.Transaction, error)
	Delete(ctx context.Context, trxId int64) error
	Find(ctx context.Context, filter model.Transaction) ([]model.Transaction, error)
	Update(ctx context.Context, trx model.Transaction) (model.Transaction, error)
	DeleteByAccountId(ctx context.Context, accountId uuid.UUID) error
}

type trxRepository struct {
	DB *gorm.DB
}

func NewTrxRepository(DB *gorm.DB) TrxRepository {
	return &trxRepository{DB: DB}
}

func (tr *trxRepository) FindAll(
	ctx context.Context,
	trxType model.TrxTypeEnum,
	startDate time.Time,
	endDate time.Time,
	accountId uuid.UUID,
	categoryId uuid.UUID,
	pagination util.Pagination,
) ([]model.Transaction, int64, error) {
	var trxs []model.Transaction
	var count int64

	query := tr.DB.WithContext(ctx).Table("transactions")

	if accountId != uuid.Nil {
		query = query.Where("account_id = ?", accountId)
	}

	if categoryId != uuid.Nil {
		query = query.Where("category_id = ?", categoryId)
	}

	query = query.Where("trx_type = ?", trxType)
	query = query.Where("trx_date >= ? AND trx_date <= ?", startDate, endDate)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	offset := (pagination.Page - 1) * pagination.PageSize

	err := query.Limit(int(pagination.PageSize)).Offset(int(offset)).Order("trx_date DESC").Find(&trxs).Error

	return trxs, count, err
}

func (tr *trxRepository) Create(ctx context.Context, trx model.Transaction) (model.Transaction, error) {
	result := tr.DB.WithContext(ctx).Create(&trx)
	return trx, result.Error
}

func (tr *trxRepository) Delete(ctx context.Context, trxId int64) error {
	return tr.DB.WithContext(ctx).Delete(&model.Transaction{}, trxId).Error
}

func (tr *trxRepository) Find(ctx context.Context, filter model.Transaction) ([]model.Transaction, error) {
	var trxs []model.Transaction
	if err := tr.DB.WithContext(ctx).Where(filter).Find(&trxs).Error; err != nil {
		return nil, err
	}

	return trxs, nil
}

func (tr *trxRepository) Update(ctx context.Context, trx model.Transaction) (model.Transaction, error) {
	result := tr.DB.WithContext(ctx).Model(&trx).Updates(trx)
	return trx, result.Error
}

func (tr *trxRepository) DeleteByAccountId(ctx context.Context, accountId uuid.UUID) error {
	return tr.DB.WithContext(ctx).Delete(&model.Transaction{}, "account_id = ?", accountId).Error
}
