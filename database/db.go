package database

import (
	"fiscus/model"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB = nil

func initDb() {
	DB_USERNAME := os.Getenv("DB_USERNAME")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_NAME := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		DB_USERNAME,
		DB_PASSWORD,
		DB_HOST,
		DB_PORT,
		DB_NAME,
	)
	// fmt.Println("MASUK")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to DB: " + err.Error())
	}
}

func GetDB() *gorm.DB {
	if DB == nil {

		initDb()
	}

	return DB
}

func MigrateTables() {
	pg := GetDB()
	pg.AutoMigrate(&model.Account{})
	pg.AutoMigrate(&model.Transaction{})
	pg.AutoMigrate(&model.Category{})
	pg.AutoMigrate(&model.RoutineExpense{})
	pg.AutoMigrate(&model.RoutinePayment{})
}
