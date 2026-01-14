package handler

import (
	"fiscus/model"
	"fiscus/service"
	"fiscus/util"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TrxHandler struct {
	ts service.TrxService
}

func NewTrxHandler(ts service.TrxService) *TrxHandler {
	return &TrxHandler{ts: ts}
}

func (th *TrxHandler) GetTransactions(ctx *gin.Context) {
	page := ctx.Query("pagenum")
	var pageNum int64
	if page == "" {
		pageNum = 1
	} else {
		var err error
		pageNum, err = strconv.ParseInt(page, 10, 64)
		if err != nil {
			fmt.Println("Failed to process request: " + err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": "false",
				"message": "Invalid Page Number",
			})
			return
		}
	}

	pageSize := ctx.Query("pagesize")
	var pageSizeNum int64
	if pageSize == "" {
		pageSizeNum = 20
	} else {
		var err error
		pageSizeNum, err = strconv.ParseInt(pageSize, 10, 64)
		if err != nil {
			fmt.Println("Failed to process request: " + err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": "false",
				"message": "Invalid Page Number",
			})
			return
		}
	}

	type Req struct {
		AccountNumber string
		Category      string
		TrxType       model.TrxTypeEnum
		StartDate     string
		EndDate       string
	}

	var req Req
	if err := ctx.ShouldBindQuery(&req); err != nil {
		fmt.Println("Failed to process request: " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Invalid Request.",
		})
		return
	}

	fmt.Println("Account Number: " + req.AccountNumber)
	startDate, err := util.StringToTime(req.StartDate)
	if err != nil {
		fmt.Println("Failed to process request: " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Invalid Start Date.",
		})
		return
	}

	endDate, err := util.StringToTime(req.EndDate)
	if err != nil {
		fmt.Println("Failed to process request: " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Invalid Start Date.",
		})
		return
	}

	if endDate.Before(startDate) {
		fmt.Println("Failed to process request: Invalid date range.")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Invalid Date Range.",
		})
		return
	}

	trxs, count, err := th.ts.FindAll(ctx, req.TrxType, startDate, endDate, req.AccountNumber, req.Category, pageNum, pageSizeNum)
	if err != nil {
		fmt.Println("Failed to process request: " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Failed to Fetch Transactions.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":      "true",
		"transactions": trxs,
		"count":        count,
	})
}

func (th *TrxHandler) CreateTransaction(ctx *gin.Context) {
	type Req struct {
		AccountNumber string
		Category      string
		TrxDate       string
		TrxType       model.TrxTypeEnum
		Amount        float64
	}

	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println("Failed to process request: " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Invalid request.",
		})
		return
	}

	accountNumber := req.AccountNumber
	if accountNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Account Number Cannot be Empty.",
		})
		return
	}

	cat := req.Category
	if cat == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Category Cannot be Empty.",
		})
		return
	}
	fmt.Println("Category: " + cat)

	trxDateParsed, err := util.StringToTime(req.TrxDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Invalid Date.",
		})
	}

	trx := model.Transaction{
		TrxDate: trxDateParsed,
		TrxType: req.TrxType,
		Amount:  req.Amount,
	}

	newTrx, err := th.ts.Create(ctx, trx, accountNumber, cat)
	if err != nil {
		fmt.Println("Error processing request : " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": "false",
			"message": "Failed to Insert Transaction.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":      "true",
		"transactions": newTrx,
	})
}
