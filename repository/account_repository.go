package repository

import (
	"context"
	"fiscus/model"
	"fmt"

	"gorm.io/gorm"
)

type AccountRepository interface {
	FindAll(context context.Context) ([]model.Account, error)
	FindByFilter(context context.Context, filter model.Account) ([]model.Account, error)
	Create(context context.Context, account model.Account) (model.Account, error)
	Update(ctx context.Context, account model.Account) (model.Account, error)
}

type accountRepository struct {
	DB *gorm.DB
}

func NewAccountRepository(DB *gorm.DB) AccountRepository {
	return &accountRepository{DB: DB}
}

func (ar *accountRepository) FindAll(context context.Context) ([]model.Account, error) {
	fmt.Println("Processing query..")
	var accounts []model.Account
	err := ar.DB.WithContext(context).Omit("Transactions").Find(&accounts).Error
	fmt.Println("Finished Processing Query..")
	return accounts, err
}

func (ar *accountRepository) FindByFilter(context context.Context, filter model.Account) ([]model.Account, error) {
	var accounts []model.Account
	err := ar.DB.WithContext(context).Where(filter).Find(&accounts).Error
	return accounts, err
}

func (ar *accountRepository) Create(context context.Context, account model.Account) (model.Account, error) {
	result := ar.DB.WithContext(context).Create(&account)
	return account, result.Error
}

func (ar *accountRepository) Update(ctx context.Context, account model.Account) (model.Account, error) {
	result := ar.DB.WithContext(ctx).Model(&account).Updates(account)
	return account, result.Error
}
