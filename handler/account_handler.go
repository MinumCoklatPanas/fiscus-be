package handler

import (
	"fiscus/model"
	"fiscus/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AccountHandler struct {
	as service.AccountService
	ts service.TrxService
}

func NewAccountHandler(as service.AccountService, ts service.TrxService) *AccountHandler {
	return &AccountHandler{as: as, ts: ts}
}

func (ah *AccountHandler) GetAllAccounts(context *gin.Context) {
	fmt.Println("Processing GetAllAcounts..")
	accounts, err := ah.as.GetAllAccounts(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"success": "false", "message": "failed to fetch account. error: " + err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"succes": "true", "accounts": accounts})
}

func (ah *AccountHandler) GetAccount(context *gin.Context) {
	var account model.Account
	err := context.ShouldBindJSON(&account)
	if err != nil {
		fmt.Println("Invalid Request: " + err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"success": "false", "message": "invalid request."})
		return
	}
	// filter := util.StructToMap(account)

	var accounts []model.Account
	accounts, err = ah.as.GetAccount(context, account)
	if err != nil {
		fmt.Println("Failed to fetch. Err : " + err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"success": "false", "message": "failed to fetch."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": "true", "accounts": accounts})
}

func (ah *AccountHandler) CreateAccount(context *gin.Context) {
	var account model.Account
	err := context.ShouldBindJSON(&account)
	if err != nil {
		fmt.Println("Invalid request: " + err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Invalid Request.",
		})
		return
	}

	var newAccount model.Account
	newAccount, err = ah.as.CreateAccount(context, account)
	if err != nil {
		fmt.Println("Failed to create account: " + err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"success": "false",
			"message": "Failed to Create Account.",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"success": "true",
		"account": newAccount,
	})
}

func (ah *AccountHandler) UpdateAccount(ctx *gin.Context) {
	accountId, err := uuid.Parse(ctx.Param("accountId"))
	if err != nil {
		fmt.Println("Failed to parse id: " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Invallid Id.",
		})
		return
	}

	var account model.Account
	if err := ctx.ShouldBindJSON(&account); err != nil {
		fmt.Println("Error Processing Payload: " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Error Processing Payload.",
		})
		return
	}

	account.Id = accountId
	newAccount, err := ah.as.UpdateAccount(ctx, account, accountId)
	if err != nil {
		fmt.Println("Error Updating Account: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": "false",
			"message": "Error Updating Account.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": "true",
		"account": newAccount,
	})
}

func (ah *AccountHandler) DeleteAccount(ctx *gin.Context) {
	accountId, err := uuid.Parse(ctx.Param("accountId"))
	if err != nil {
		fmt.Println("Failed to parse id: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": "false",
			"message": "Invalid Account Id.",
		})
		return
	}

	err = ah.ts.DeleteByAccountId(ctx, accountId)
	if err != nil {
		fmt.Println("Failed to delte transactions: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": "false",
			"message": "Failed to Delete Transactions.",
		})
		return
	}

	err = ah.as.DeleteAccount(ctx, accountId)
	if err != nil {
		fmt.Println("Failed to delete account: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": "false",
			"message": "Failed to Delete Account.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": "true",
		"message": "Successfully Deleted Account",
	})
}
