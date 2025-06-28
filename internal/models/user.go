package models

import (
	"time"
)

type User struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Email     string `gorm:"type:varchar(100);uniqueIndex"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Orgs      []Organization `gorm:"many2many:user_organizations;"`
}

type Organization struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Users     []User `gorm:"many2many:user_organizations;"`
}

type UserOrganization struct {
	UserID         string `gorm:"primaryKey"`
	OrganizationID string `gorm:"primaryKey"`
	IsCurrent      bool
}
