package database

import (
	"log"

	"gitag.ir/cookthepot/services/vault/config"
	"gitag.ir/cookthepot/services/vault/database/migrations"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if config.AppConfig.Environment != "development" {
		log.Fatal("you are not in development mode. migration wont start !!!")
	}

	// migrate the dependent many2many before the main table
	migrations.M221Category(db)
	migrations.M223Attachment(db)
	migrations.M231Role(db)
	migrations.M222User(db)
	migrations.M228Verification(db)
	migrations.M281Token(db)
}
