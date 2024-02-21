package database

import (
	"log"

	"github.com/salimkun/vending-machine-api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Instance *gorm.DB
var dbError error

func Connect() {
	Instance, dbError = gorm.Open(sqlite.Open("vendingMachine.db"), &gorm.Config{})
	if dbError != nil {
		log.Fatal(dbError)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Database!")
}

func Migrate() {
	Instance.AutoMigrate(&models.User{})
	Instance.AutoMigrate(&models.Product{})
	log.Println("Database Migration Completed!")
}
