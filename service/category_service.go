package service

import (
	"context"
	"errors"
	"fiscus/model"
	"fiscus/repository"
)

type CategoryService interface {
	GetAllCategories(ctx context.Context) ([]model.Category, error)
	GetCategoryByName(ctx context.Context, name string) (model.Category, error)
	CreateCategory(ctx context.Context, category model.Category) (model.Category, error)
}

type categoryService struct {
	cr repository.CategoryRepository
}

func NewCategoryService(cr repository.CategoryRepository) CategoryService {
	return &categoryService{cr: cr}
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
