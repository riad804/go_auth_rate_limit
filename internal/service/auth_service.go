package service

import (
	"errors"
	"time"

	"github.com/riad804/go_auth/internal/repository"
	"github.com/riad804/go_auth/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  repository.UserRepository
	tokenRepo *repository.TokenRepository
	jwtUtil   *utils.JWTWrapper
}

func NewAuthService(ur repository.UserRepository, tr *repository.TokenRepository, jwt *utils.JWTWrapper) *AuthService {
	return &AuthService{userRepo: ur, tokenRepo: tr, jwtUtil: jwt}
}

func (s *AuthService) Login(email, password string) (accessToken, refreshToken string, err error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	currentOrg, err := s.userRepo.GetCurrentOrg(user.ID)
	if err != nil {
		return "", "", errors.New("organization error")
	}

	accessToken, err = s.jwtUtil.GenerateToken(user.ID, currentOrg.ID)
	if err != nil {
		return "", "", errors.New("token generation failed")
	}

	refreshToken, err = s.jwtUtil.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", errors.New("refresh token failed")
	}

	if err := s.tokenRepo.StoreRefreshToken(refreshToken, user.ID, 7*24*time.Hour); err != nil {
		return "", "", errors.New("token storage failed")
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(refreshToken string) (string, string, error) {
	userID, err := s.tokenRepo.GetRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	currentOrg, err := s.userRepo.GetCurrentOrg(user.ID)
	if err != nil {
		return "", "", errors.New("organization error")
	}

	accessToken, err := s.jwtUtil.GenerateToken(user.ID, currentOrg.ID)
	if err != nil {
		return "", "", errors.New("token generation failed")
	}

	newRefreshToken, err := s.jwtUtil.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", errors.New("refresh token failed")
	}

	if err := s.tokenRepo.StoreRefreshToken(newRefreshToken, user.ID, 7*24*time.Hour); err != nil {
		return "", "", errors.New("token storage failed")
	}

	if err := s.tokenRepo.DeleteRefreshToken(refreshToken); err != nil {
		return "", "", errors.New("token deletion failed")
	}

	return accessToken, newRefreshToken, nil
}

func (s *AuthService) SwitchOrg(userID, orgID string) (string, error) {
	if err := s.userRepo.SwitchCurrentOrg(userID, orgID); err != nil {
		return "", errors.New("failed to switch organization")
	}

	currentOrg, err := s.userRepo.GetCurrentOrg(userID)
	if err != nil {
		return "", errors.New("failed to get current organization")
	}

	accessToken, err := s.jwtUtil.GenerateToken(userID, currentOrg.ID)
	if err != nil {
		return "", errors.New("token generation failed")
	}

	return accessToken, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.tokenRepo.DeleteRefreshToken(refreshToken)
}
