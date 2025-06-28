package database

import (
	"fmt"

	"github.com/riad804/go_auth/internal/config"
	"github.com/riad804/go_auth/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Organization{}, &models.UserOrganization{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
