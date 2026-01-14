package main

import (
	"fiscus/handler"
	"fiscus/repository"
	"fiscus/service"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAccountRoutes(db *gorm.DB, server *gin.RouterGroup) {
	ar := repository.NewAccountRepository(db)
	as := service.NewAccountService(ar)

	cr := repository.NewCategoryRepository(db)
	cs := service.NewCategoryService(cr)

	tr := repository.NewTrxRepository(db)
	ts := service.NewTrxService(tr, as, cs)

	ah := handler.NewAccountHandler(as, ts)

	fmt.Println("Registering Account Routes..")
	server.GET("/accounts", ah.GetAllAccounts)
	server.GET("/account", ah.GetAccount)
	server.POST("/account", ah.CreateAccount)
	server.PATCH("/account/:accountId", ah.UpdateAccount)
	server.DELETE("/account/:accountId", ah.DeleteAccount)
}

func RegisterCategoryRoutes(db *gorm.DB, server *gin.RouterGroup) {
	cr := repository.NewCategoryRepository(db)
	cs := service.NewCategoryService(cr)
	ch := handler.NewCategoryHandler(cs)

	fmt.Println("Registering Category Routes..")
	server.GET("/categories", ch.GetAllCategories)
	server.POST("/category", ch.CreateCategory)
}

func RegisterTrxRoutes(db *gorm.DB, server *gin.RouterGroup) {
	ar := repository.NewAccountRepository(db)
	as := service.NewAccountService(ar)

	cr := repository.NewCategoryRepository(db)
	cs := service.NewCategoryService(cr)

	tr := repository.NewTrxRepository(db)
	ts := service.NewTrxService(tr, as, cs)

	th := handler.NewTrxHandler(ts)

	fmt.Println("Registering Trx Routes..")
	server.GET("/transactions", th.GetTransactions)
	server.POST("/transaction", th.CreateTransaction)
	server.DELETE("/transaction/:id", th.DeleteTransaction)
}

func RegisterAuthRoutes(server *gin.Engine) {
	ah := handler.NewAuthHandler()
	fmt.Println("Registering Auth Routes..")
	server.POST("/login", ah.Login)
}
