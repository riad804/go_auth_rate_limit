package test

import (
	"testing"
	"time"

	"github.com/riad804/go_auth/internal/utils"
	"github.com/stretchr/testify/assert"
)

type dummyConfig struct{}

func (d *dummyConfig) JWTSecret() string              { return "testsecret" }
func (d *dummyConfig) AccessTokenExp() time.Duration  { return time.Minute }
func (d *dummyConfig) RefreshTokenExp() time.Duration { return time.Hour }

func TestToken(t *testing.T) {
	jwt := utils.JWTWrapper{
		SecretKey:       "testsecret",
		AccessTokenExp:  time.Minute,
		RefreshTokenExp: time.Hour,
	}

	token, err := jwt.GenerateToken("user1", "org1")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := jwt.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "user1", claims.UserID)
	assert.Equal(t, "org1", claims.OrgID)
}

func TestTokenInvalid(t *testing.T) {
	jwt := utils.JWTWrapper{
		SecretKey:       "testsecret",
		AccessTokenExp:  time.Minute,
		RefreshTokenExp: time.Hour,
	}

	_, err := jwt.ValidateToken("invalid.token")
	assert.Error(t, err)
}
