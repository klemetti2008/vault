package seeders

import (
	"gitag.ir/cookthepot/services/vault/models"
	"github.com/mhosseintaher/kit/dtp"
	"github.com/mhosseintaher/kit/faker"
	"gorm.io/gorm"

	"log"
)

func UserSeeder(db *gorm.DB) {
	users := []models.User{
		{
			Bio:             "I am a highly motivated and experienced Software Engineer with a passion for solving complex problems through innovative software solutions.",
			Title:           "Managing Director",
			Name:            "Matthew",
			Phone:           "720-931-1196",
			Grade:           1000,
			LastName:        "L. Gonzales",
			Nickname:        "Mr.Matt",
			AvatarUrl:       "avatars/5c9abba1-a285-4f4c-83e2-4db4f1b3fbb5",
			Password:        "$2a$14$ALuVvY5tXovI9WFyxed/eORguxkDhWaEgHARuOMdekOloWFL2zuMO", // password
			CountryCode:     "US",
			PhoneVerifiedAt: faker.SQLNow(),
			IDCode: dtp.NullString{
				String: "008-84-0497",
				Valid:  true,
			},
			Username: dtp.NullString{
				String: "mattgoz",
				Valid:  true,
			},
			Email: dtp.NullString{
				String: "matt.gonzales@gmail.com",
				Valid:  true,
			},
			EmailVerifiedAt:    faker.SQLNow(),
			ProfileCompletedAt: faker.SQLNow(),
			MadeOfficialAt:     faker.SQLNow(),
			SuspendedAt:        dtp.NullTime{},
			Roles: []*models.Role{
				{ID: 1},
				{ID: 2},
				{ID: 3},
			},
		},
		{
			Name:            "Theresa",
			Phone:           "952-525-1190",
			LastName:        "Sewell",
			CountryCode:     "US",
			Password:        "$2a$14$ALuVvY5tXovI9WFyxed/eORguxkDhWaEgHARuOMdekOloWFL2zuMO", // password
			PhoneVerifiedAt: faker.SQLNow(),
			Username: dtp.NullString{
				String: "698-12-9881",
				Valid:  true,
			},
			Email: dtp.NullString{
				String: "theresa.sewell@gmail.com",
				Valid:  true,
			},
			EmailVerifiedAt: faker.SQLNow(),
			Roles: []*models.Role{
				{ID: 1},
				{ID: 2},
			},
		},
	}

	err := db.Create(&users).Error
	if err != nil {
		log.Fatal(err)
	}
}
