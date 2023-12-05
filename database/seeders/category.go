package seeders

import (
	"log"

	"gitag.ir/cookthepot/services/vault/models"
	"github.com/mhosseintaher/kit/dtp"
	"gorm.io/gorm"
)

func CategorySeeder(db *gorm.DB) {
	categories := []models.Category{
		{
			// ID: 1
			ParentID: dtp.NullInt64{},
			Title:    "test1",
			UserID:   1,
		},
		{
			// ID : 2
			ParentID: dtp.NullInt64{
				Int64: 1,
				Valid: true,
			},
			Title:  "test2",
			UserID: 1,
		},
		{
			// ID: 3
			ParentID: dtp.NullInt64{},
			Title:    "test3",
			UserID:   1,
		},
	}

	err := db.Create(&categories).Error
	if err != nil {
		log.Fatal(err)
	}
}
