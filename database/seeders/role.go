package seeders

import (
	"gitag.ir/cookthepot/services/vault/models"
	"gitag.ir/cookthepot/services/vault/services/role"
	"gorm.io/gorm"
	"log"
)

func RoleSeeder(db *gorm.DB) {

	var roles []models.Role
	for _, v := range role.Roles {
		roles = append(
			roles, models.Role{
				Title: v,
			},
		)
	}
	err := db.Create(&roles).Error
	if err != nil {
		log.Fatal(err)
	}
}
