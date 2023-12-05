package database

import (
	"gitag.ir/cookthepot/services/vault/models"
	"gorm.io/gorm"

	"log"
)

func DropJoinTables(db *gorm.DB) {
	// TODO: handle errors
	err := db.Migrator().DropTable("user_craft")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrator().DropTable("user_permission")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrator().DropTable("user_role")
	if err != nil {
		log.Fatal(err)
	}
}

func DropTables(db *gorm.DB) {
	err := db.Migrator().DropTable(&models.User{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrator().DropTable(&models.Attachment{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrator().DropTable(&models.Category{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrator().DropTable(&models.Verification{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrator().DropTable(&models.Role{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().DropTable(&models.Token{})
	if err != nil {
		log.Fatal(err)
	}
}

func DropAll(db *gorm.DB) {
	DropJoinTables(db.Debug())
	DropTables(db.Debug())
}
