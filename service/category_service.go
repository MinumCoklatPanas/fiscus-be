package service

import (
	"context"
	"errors"
	"fiscus/dto"
	"fiscus/model"
	"fiscus/repository"
	"time"

	"github.com/google/uuid"
)

type CategoryService interface {
	GetAllCategories(ctx context.Context) ([]model.Category, error)
	GetCategoryByName(ctx context.Context, name string) (model.Category, error)
	CreateCategory(ctx context.Context, category model.Category) (model.Category, error)
	GetSumPerCategory(
		ctx context.Context,
		trxType model.TrxTypeEnum,
		startDate time.Time,
		endDate time.Time,
		accountNumber string,
		categoryName string,
	) ([]dto.TransactionPerCategory, error)
}

type categoryService struct {
	cr repository.CategoryRepository
	as AccountService
}

func NewCategoryService(cr repository.CategoryRepository, as AccountService) CategoryService {
	return &categoryService{cr: cr, as: as}
}

func (cs *categoryService) GetAllCategories(ctx context.Context) ([]model.Category, error) {
	return cs.cr.FindAll(ctx)
}

func (cs *categoryService) GetCategoryByName(ctx context.Context, name string) (model.Category, error) {
	categories, err := cs.cr.FindByFilter(ctx, model.Category{Name: name})
	if err != nil {
		return model.Category{}, err
	}
	if len(categories) == 0 {
		return model.Category{}, errors.New("Category Not Found.")
	}
	return categories[0], err
}

func (cs *categoryService) CreateCategory(ctx context.Context, category model.Category) (model.Category, error) {
	return cs.cr.Create(ctx, category)
}

func (cs *categoryService) GetSumPerCategory(
	ctx context.Context,
	trxType model.TrxTypeEnum,
	startDate time.Time,
	endDate time.Time,
	accountNumber string,
	categoryName string,
) ([]dto.TransactionPerCategory, error) {
	var accountId uuid.UUID
	if accountNumber != "" {
		account, err := cs.as.GetAccountByNumber(ctx, accountNumber)
		if err != nil {
			return nil, err
		}
		accountId = account.Id
	}

	var categoryId uuid.UUID
	if categoryName != "" {
		cat, err := cs.GetCategoryByName(ctx, categoryName)
		if err != nil {
			return nil, err
		}
		categoryId = cat.Id
	}

	return cs.cr.SumPerTrx(ctx, trxType, startDate, endDate, accountId, categoryId)
}
