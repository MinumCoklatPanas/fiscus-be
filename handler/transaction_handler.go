package handler

import (
	"fiscus/model"
	"fiscus/service"
	"fiscus/util"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TrxHandler struct {
	ts service.TrxService
	cs service.CategoryService
}

func NewTrxHandler(ts service.TrxService, cs service.CategoryService) *TrxHandler {
	return &TrxHandler{ts: ts, cs: cs}
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

	sumPerCategory, err := th.cs.GetSumPerCategory(ctx, req.TrxType, startDate, endDate, req.AccountNumber, req.Category)
	if err != nil {
		fmt.Println("Failed to process request: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": "false",
			"message": "Failed to get statistics.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":      "true",
		"transactions": trxs,
		"statistics":   sumPerCategory,
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
		Notes         *string
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
		Notes:   req.Notes,
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

func (th *TrxHandler) DeleteTransaction(ctx *gin.Context) {
	trxId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		fmt.Println("Error Parsing Id: " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Invalid Transaction Id.",
		})
		return
	}

	err = th.ts.Delete(ctx, int64(trxId))
	if err != nil {
		fmt.Println("Failed to delete transaction: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": "false",
			"message": "Failed to Delete Transaction.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": "true",
		"message": "Successfully Deleted Transction.",
	})
}

func (th *TrxHandler) UpdateTransaction(ctx *gin.Context) {
	trxId := ctx.Param("id")

	type Req struct {
		Amount     float64
		TrxDate    string
		CategoryId uuid.UUID
		Notes      *string
		TrxType    model.TrxTypeEnum
	}
	var req Req

	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println("Failed Processing Payload : " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Invalid Request.",
		})
		return
	}

	trxDateParsed := time.Time{}
	var err error
	if req.TrxDate != "" {
		trxDateParsed, err = util.StringToTime(req.TrxDate)
		if err != nil {
			fmt.Println("Failed to parse trx date: " + err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": "false",
				"message": "Invalid Transaction Date.",
			})
			return
		}
	}

	idNum, err := strconv.ParseInt(trxId, 10, 64)
	if err != nil {
		fmt.Println("Invalid Transaction Id : " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Invalid Transaction Id.",
		})
		return
	}

	trx := model.Transaction{
		Id:         idNum,
		Amount:     req.Amount,
		TrxDate:    trxDateParsed,
		CategoryId: req.CategoryId,
		Notes:      req.Notes,
		TrxType:    req.TrxType,
	}

	newTrx, err := th.ts.Update(ctx, trx, idNum)
	if err != nil {
		fmt.Println("Failed to Update Transaction : " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": "false",
			"message": "Failed to Update Transaction.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":     "true",
		"transaction": newTrx,
	})
}
