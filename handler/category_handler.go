package handler

import (
	"fiscus/model"
	"fiscus/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	cs service.CategoryService
}

func NewCategoryHandler(cs service.CategoryService) *CategoryHandler {
	return &CategoryHandler{cs: cs}
}

func (ch *CategoryHandler) GetAllCategories(ctx *gin.Context) {
	categories, err := ch.cs.GetAllCategories(ctx)
	if err != nil {
		fmt.Println("Failed to process request: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"succes":  "false",
			"message": "Failed to fetch categories.",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success":    "true",
		"categories": categories,
	})
}

func (ch *CategoryHandler) CreateCategory(ctx *gin.Context) {
	var category model.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		fmt.Println("Failed to process request: " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Invalid category.",
		})
		return
	}

	var newCategory model.Category
	newCategory, err := ch.cs.CreateCategory(ctx, category)
	if err != nil {
		fmt.Println("Failed to process request: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": "false",
			"message": "Failed to insert category.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":  "true",
		"category": newCategory,
	})
}
