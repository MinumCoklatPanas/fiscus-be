package service

import (
	"context"
	"errors"
	"fiscus/model"
	"fiscus/repository"
	"fiscus/util"
	"fmt"
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
	Delete(ctx context.Context, trxId int64) error
	Update(ctx context.Context, newTrx model.Transaction, trxId int64) (model.Transaction, error)
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

func (ts *trxService) Delete(ctx context.Context, trxId int64) error {
	var trxs []model.Transaction
	trxs, err := ts.tr.Find(ctx, model.Transaction{Id: trxId})
	if err != nil {
		return err
	}

	if len(trxs) == 0 {
		return errors.New("Cannot Find Transaction.")
	}

	curTrx := trxs[0]
	delta := curTrx.Amount
	if curTrx.TrxType == "income" {
		delta *= -1
	}

	account, err := ts.as.GetAccount(ctx, model.Account{Id: curTrx.AccountId})
	if err != nil {
		return err
	}

	if len(account) == 0 {
		return errors.New("Cannot Find Account.")
	}

	if len(account) > 1 {
		return errors.New("Invalid Account.")
	}

	_, err = ts.as.UpdateAccountBalance(ctx, account[0], delta)
	if err != nil {
		return err
	}

	return ts.tr.Delete(ctx, trxId)
}

func (ts *trxService) Update(ctx context.Context, newTrx model.Transaction, trxId int64) (model.Transaction, error) {
	existingTrx, err := ts.tr.Find(ctx, model.Transaction{Id: trxId})
	if err != nil {
		return model.Transaction{}, err
	}

	if len(existingTrx) == 0 {
		return model.Transaction{}, errors.New("Multiple Transaction Detected.")
	}

	if len(existingTrx) > 1 {
		return model.Transaction{}, errors.New("Duplicate Transaction Id Found.")
	}

	if newTrx.Amount == 0 {
		newTrx.Amount = existingTrx[0].Amount
	} else {
		var delta float64
		if existingTrx[0].TrxType == "expense" {
			delta = existingTrx[0].Amount - newTrx.Amount
		} else {
			delta = newTrx.Amount - existingTrx[0].Amount
		}

		account, err := ts.as.GetAccount(ctx, model.Account{Id: existingTrx[0].AccountId})
		if err != nil {
			return model.Transaction{}, err
		}

		_, err = ts.as.UpdateAccountBalance(ctx, account[0], delta)
		if err != nil {
			return model.Transaction{}, err
		}
	}

	if newTrx.TrxDate.IsZero() {
		newTrx.TrxDate = existingTrx[0].TrxDate
	}

	if newTrx.CategoryId != uuid.Nil {
		newTrx.CategoryId = existingTrx[0].CategoryId
	}

	if newTrx.Notes == nil {
		newTrx.Notes = existingTrx[0].Notes
	}

	if newTrx.TrxType == "" {
		newTrx.TrxType = existingTrx[0].TrxType
	}

	fmt.Println("Updating Transaction..")
	return ts.tr.Update(ctx, newTrx)
}

func (ts *trxService) DeleteByAccountId(ctx context.Context, accountId uuid.UUID) error {
	return ts.tr.DeleteByAccountId(ctx, accountId)
}
