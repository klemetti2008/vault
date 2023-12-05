package database

import (
	"gitag.ir/cookthepot/services/vault/database/seeders"
	"gorm.io/gorm"
)

func SeedAllTables(db *gorm.DB) {
	seeders.RoleSeeder(db)
	seeders.UserSeeder(db)
}
