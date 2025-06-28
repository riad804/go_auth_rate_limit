package main

import (
	"fmt"

	"github.com/riad804/go_auth/internal/models"
	"github.com/riad804/go_auth/pkg/database"

	"github.com/riad804/go_auth/internal/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()
	db, err := database.NewMySQLDB(cfg)
	if err != nil {
		panic("failed to connect database")
	}

	seedData(db)
}

func seedData(db *gorm.DB) {
	orgs := []models.Organization{
		{ID: "org1", Name: "Tenbyte"},
		{ID: "org2", Name: "OpenResty"},
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	users := []models.User{
		{
			ID:       "user1",
			Name:     "Sohel",
			Email:    "sohel@tenbyte.com",
			Password: string(hashedPassword),
			Orgs:     []models.Organization{orgs[0], orgs[1]},
		},
		{
			ID:       "user2",
			Name:     "Jane",
			Email:    "jane@openresty.com",
			Password: string(hashedPassword),
			Orgs:     []models.Organization{orgs[0], orgs[1]},
		},
		{
			ID:       "user3",
			Name:     "Riad",
			Email:    "riad@openresty.com",
			Password: string(hashedPassword),
			Orgs:     []models.Organization{orgs[0], orgs[1]},
		},
		{
			ID:       "user4",
			Name:     "Tanvir",
			Email:    "tanvir@tenbyte.com",
			Password: string(hashedPassword),
			Orgs:     []models.Organization{orgs[0], orgs[1]},
		},
	}

	for _, org := range orgs {
		if err := db.Create(&org).Error; err != nil {
			fmt.Printf("Error seeding org %s: %v\n", org.Name, err)
		}
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			fmt.Printf("Error seeding user %s: %v\n", user.Name, err)
		}

		db.Create(&models.UserOrganization{
			UserID:         user.ID,
			OrganizationID: user.Orgs[0].ID,
			IsCurrent:      true,
		})
	}
}
