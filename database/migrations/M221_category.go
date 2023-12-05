package migrations

import (
	"gitag.ir/cookthepot/services/vault/models"
	"gorm.io/gorm"

	"log"
)

func M221Category(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&models.Category{})

	if err != nil {
		log.Fatal(err)
	}
}
