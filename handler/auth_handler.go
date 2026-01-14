package handler

import (
	"fiscus/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (ah *AuthHandler) Login(ctx *gin.Context) {
	type Req struct {
		Username string
		Password string
	}

	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println("Failed to process request: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": "false",
			"message": "Invalid login.",
		})
		return
	}
	// fmt.Println("[DEBUG] " + req.Password + " " + req.Username)

	isVerified := util.VerifyUser(req.Username, req.Password)
	if !isVerified {
		fmt.Println("Failed to process request: invalid credentials.")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": "Invalid Credentials.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": "true",
	})
}
