package service

import (
	"github.com/riad804/go_auth/internal/models"
	"github.com/riad804/go_auth/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserWithOrgs(userID string) (*models.User, []*models.Organization, error) {
	return s.repo.GetUserWithOrgs(userID)
}

func (s *UserService) GetCurrentOrg(userID string) (*models.Organization, error) {
	return s.repo.GetCurrentOrg(userID)
}
