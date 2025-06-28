package repository

import (
	"github.com/riad804/go_auth/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	GetCurrentOrg(userID string) (*models.Organization, error)
	GetUserWithOrgs(userID string) (*models.User, []*models.Organization, error)
	SwitchCurrentOrg(userID, orgID string) error
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetCurrentOrg(userID string) (*models.Organization, error) {
	var userOrg models.UserOrganization
	if err := r.db.Where("user_id = ? AND is_current = ?", userID, true).First(&userOrg).Error; err != nil {
		return nil, err
	}
	var org models.Organization
	if err := r.db.Where("id = ?", userOrg.OrganizationID).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *GormUserRepository) GetUserWithOrgs(userID string) (*models.User, []*models.Organization, error) {
	var user models.User
	if err := r.db.Preload("Orgs").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, nil, err
	}
	orgs := make([]*models.Organization, len(user.Orgs))
	for i := range user.Orgs {
		orgs[i] = &user.Orgs[i]
	}
	return &user, orgs, nil
}

func (r *GormUserRepository) SwitchCurrentOrg(userID, orgID string) error {
	// Set all user's orgs to IsCurrent = false
	if err := r.db.Model(&models.UserOrganization{}).
		Where("user_id = ?", userID).
		Update("is_current", false).Error; err != nil {
		return err
	}
	// Set the selected org to IsCurrent = true
	if err := r.db.Model(&models.UserOrganization{}).
		Where("user_id = ? AND organization_id = ?", userID, orgID).
		Update("is_current", true).Error; err != nil {
		return err
	}
	return nil
}
