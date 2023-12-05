package migrations

import (
	"gitag.ir/cookthepot/services/vault/models"
	"gorm.io/gorm"
	"log"
)

func M222User(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
	}
}
