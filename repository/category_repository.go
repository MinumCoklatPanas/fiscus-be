package repository

import (
	"context"
	"fiscus/model"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll(ctx context.Context) ([]model.Category, error)
	FindByFilter(ctx context.Context, filter model.Category) ([]model.Category, error)
	Create(ctx context.Context, category model.Category) (model.Category, error)
}

type categoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(DB *gorm.DB) CategoryRepository {
	return &categoryRepository{DB: DB}
}

func (cr *categoryRepository) FindAll(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	err := cr.DB.WithContext(ctx).Find(&categories).Error
	return categories, err
}

func (cr *categoryRepository) FindByFilter(ctx context.Context, filter model.Category) ([]model.Category, error) {
	var categories []model.Category
	err := cr.DB.WithContext(ctx).Where(filter).Find(&categories).Error
	return categories, err
}

func (cr *categoryRepository) Create(ctx context.Context, category model.Category) (model.Category, error) {
	result := cr.DB.Create(&category)
	return category, result.Error
}
