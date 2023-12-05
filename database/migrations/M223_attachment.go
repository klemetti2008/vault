package migrations

import (
	"gitag.ir/cookthepot/services/vault/models"
	"gorm.io/gorm"
	"log"
)

func M223Attachment(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&models.Attachment{})
	if err != nil {
		log.Fatal(err)
	}
}
