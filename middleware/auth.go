package middleware

import (
	"encoding/base64"
	"fiscus/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authorize(ctx *gin.Context) {
	auth := ctx.Request.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Basic ") {
		ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": "false",
			"message": "Authorization Required.",
		})
		return
	}

	payload := strings.TrimPrefix(auth, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": "false",
			"message": "Invalid Authorization",
		})
		return
	}

	splitted := strings.Split(string(decoded), ":")
	userName := splitted[0]
	password := splitted[1]

	verified := util.VerifyUser(userName, password)

	if !verified {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": "false",
			"message": "Invalid Authorization",
		})
		return
	}

	ctx.Set("username", strings.ToLower(userName))
	ctx.Next()
}
