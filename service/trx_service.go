package service

import (
	"context"
	"fiscus/model"
	"fiscus/repository"
	"fiscus/util"
	"time"

	"github.com/google/uuid"
)

type TrxService interface {
	FindAll(
		ctx context.Context,
		trxType model.TrxTypeEnum,
		startDate time.Time,
		endDate time.Time,
		accountNumber string,
		category string,
		page int64,
		pageSize int64,
	) ([]model.Transaction, int64, error)

	Create(ctx context.Context, trx model.Transaction, accountNumber string, category string) (model.Transaction, error)
	DeleteByAccountId(ctx context.Context, accountId uuid.UUID) error
}

type trxService struct {
	tr repository.TrxRepository
	as AccountService
	cs CategoryService
}

func NewTrxService(tr repository.TrxRepository, as AccountService, cs CategoryService) TrxService {
	return &trxService{
		tr: tr,
		as: as,
		cs: cs,
	}
}

func (ts *trxService) FindAll(
	ctx context.Context,
	trxType model.TrxTypeEnum,
	startDate time.Time,
	endDate time.Time,
	accountNumber string,
	category string,
	page int64,
	pageSize int64,
) ([]model.Transaction, int64, error) {

	var accountId uuid.UUID
	if accountNumber != "" {
		account, err := ts.as.GetAccountByNumber(ctx, accountNumber)
		if err != nil {
			return []model.Transaction{}, 0, err
		}
		accountId = account.Id
	}

	var categoryId uuid.UUID
	if category != "" {
		cat, err := ts.cs.GetCategoryByName(ctx, category)
		if err != nil {
			return []model.Transaction{}, 0, err
		}
		categoryId = cat.Id
	}

	pagination := util.Pagination{
		Page:     page,
		PageSize: pageSize,
	}
	pagination.Normalize()

	return ts.tr.FindAll(ctx, trxType, startDate, endDate, accountId, categoryId, pagination)
}

func (ts *trxService) Create(ctx context.Context, trx model.Transaction, accountNumber string, category string) (model.Transaction, error) {
	account, err := ts.as.GetAccountByNumber(ctx, accountNumber)
	if err != nil {
		return model.Transaction{}, err
	}
	trx.AccountId = account.Id

	delta := float64(trx.Amount)
	if trx.TrxType == "expense" {
		delta *= -1
	}

	_, err = ts.as.UpdateAccountBalance(ctx, account, delta)
	if err != nil {
		return model.Transaction{}, err
	}

	categoryObject, err := ts.cs.GetCategoryByName(ctx, category)
	if err != nil {
		return model.Transaction{}, err
	}
	trx.CategoryId = categoryObject.Id

	return ts.tr.Create(ctx, trx)
}

func (ts *trxService) DeleteByAccountId(ctx context.Context, accountId uuid.UUID) error {
	return ts.tr.DeleteByAccountId(ctx, accountId)
}
