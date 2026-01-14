package service

import (
	"context"
	"errors"
	"fiscus/model"
	"fiscus/repository"
	"fmt"

	"github.com/google/uuid"
)

type AccountService interface {
	GetAllAccounts(context context.Context) ([]model.Account, error)
	GetAccount(context context.Context, filter model.Account) ([]model.Account, error)
	GetAccountByNumber(ctx context.Context, accountNumber string) (model.Account, error)
	CreateAccount(context context.Context, account model.Account) (model.Account, error)
	UpdateAccountBalance(ctx context.Context, account model.Account, delta float64) (model.Account, error)
	UpdateAccount(ctx context.Context, account model.Account, accountId uuid.UUID) (model.Account, error)
}

type accountService struct {
	ar repository.AccountRepository
}

func NewAccountService(ar repository.AccountRepository) AccountService {
	return &accountService{ar: ar}
}

func (as *accountService) GetAllAccounts(context context.Context) ([]model.Account, error) {
	fmt.Println("Processing on service level..")
	return as.ar.FindAll(context)
}

func (as *accountService) GetAccount(context context.Context, filter model.Account) ([]model.Account, error) {
	return as.ar.FindByFilter(context, filter)
}

func (as *accountService) GetAccountByNumber(context context.Context, accountNumber string) (model.Account, error) {
	filter := model.Account{AccountNumber: accountNumber}
	var accounts []model.Account
	accounts, err := as.GetAccount(context, filter)
	if len(accounts) == 0 {
		return model.Account{}, errors.New("Account not found.")
	}
	return accounts[0], err
}

func (as *accountService) CreateAccount(context context.Context, account model.Account) (model.Account, error) {
	return as.ar.Create(context, account)
}

func (as *accountService) UpdateAccountBalance(ctx context.Context, account model.Account, delta float64) (model.Account, error) {
	account.Balance += delta
	return as.ar.Update(ctx, account)
}

func (as *accountService) UpdateAccount(ctx context.Context, newAccount model.Account, accountId uuid.UUID) (model.Account, error) {
	existingAccount, err := as.GetAccount(ctx, model.Account{Id: accountId})
	if err != nil {
		return model.Account{}, err
	}

	if newAccount.OwnerName == "" {
		newAccount.OwnerName = existingAccount[0].OwnerName
	}

	if newAccount.AccountNumber == "" {
		newAccount.AccountNumber = existingAccount[0].AccountNumber
	}

	if newAccount.AccountType == "" {
		newAccount.AccountType = existingAccount[0].AccountType
	}

	if newAccount.Balance == 0 {
		newAccount.Balance = existingAccount[0].Balance
	}

	if newAccount.ProviderName == nil {
		newAccount.ProviderName = existingAccount[0].ProviderName
	}

	if newAccount.Purpose == nil {
		newAccount.Purpose = existingAccount[0].Purpose
	}

	return as.ar.Update(ctx, newAccount)
}
