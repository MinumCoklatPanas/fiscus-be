package repository

import (
	"context"
	"fiscus/dto"
	"fiscus/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll(ctx context.Context) ([]model.Category, error)
	FindByFilter(ctx context.Context, filter model.Category) ([]model.Category, error)
	Create(ctx context.Context, category model.Category) (model.Category, error)
	SumPerTrx(
		ctx context.Context,
		trxType model.TrxTypeEnum,
		startDate time.Time,
		endDate time.Time,
		accountId uuid.UUID,
		categoryId uuid.UUID,
	) ([]dto.TransactionPerCategory, error)
}

type categoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(DB *gorm.DB) CategoryRepository {
	return &categoryRepository{DB: DB}
}

func (cr *categoryRepository) FindAll(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	err := cr.DB.WithContext(ctx).Order("name ASC").Find(&categories).Error
	return categories, err
}

func (cr *categoryRepository) FindByFilter(ctx context.Context, filter model.Category) ([]model.Category, error) {
	var categories []model.Category
	err := cr.DB.WithContext(ctx).Where(filter).Order("name ASC").Find(&categories).Error
	return categories, err
}

func (cr *categoryRepository) Create(ctx context.Context, category model.Category) (model.Category, error) {
	result := cr.DB.Create(&category)
	return category, result.Error
}

func (cr *categoryRepository) SumPerTrx(
	ctx context.Context,
	trxType model.TrxTypeEnum,
	startDate time.Time,
	endDate time.Time,
	accountId uuid.UUID,
	categoryId uuid.UUID,
) ([]dto.TransactionPerCategory, error) {
	var result []dto.TransactionPerCategory
	query := cr.DB.WithContext(ctx).Table("categories")

	if accountId != uuid.Nil {
		query = query.Where("account_id = ?", accountId)
	}

	if categoryId != uuid.Nil {
		query = query.Where("id = ?", categoryId)
	}

	query = query.Where("trx_type = ?", trxType)
	query = query.Where("trx_date >= ? AND trx_date <= ?", startDate, endDate)

	err := query.Select("categories.name as category_name, sum(transactions.amount) as total").Joins("join transactions on transactions.category_id = categories.id").Group("categories.name").Scan(&result).Error

	return result, err
}
