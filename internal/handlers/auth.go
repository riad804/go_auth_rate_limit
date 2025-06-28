package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/riad804/go_auth/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewAuthHandler(authService *service.AuthService, userService *service.UserService) *AuthHandler {
	return &AuthHandler{authService: authService, userService: userService}
}

// Login godoc
// @Summary      User login
// @Description  Authenticates user and returns access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body  object{email=string,password=string}  true  "User credentials"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	accessToken, refreshToken, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Me godoc
// @Summary      Get current user info
// @Description  Returns user info and organizations for the authenticated user
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /me [get]
// @Security     BearerAuth
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")
	user, orgs, err := h.userService.GetUserWithOrgs(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	currentOrg, err := h.userService.GetCurrentOrg(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, gin.H{
		"user": gin.H{
			"id":   user.ID,
			"name": user.Name,
		},
		"current_org": gin.H{
			"id":   currentOrg.ID,
			"name": currentOrg.Name,
		},
		"orgs": orgs,
	})
}

// SwitchOrg godoc
// @Summary      Switch organization
// @Description  Switches the current organization for the user and returns a new access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        org  body  object{org_id=string}  true  "Organization ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /orgs/switch [post]
// @Security     BearerAuth
func (h *AuthHandler) SwitchOrg(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		OrgID string `json:"org_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	accessToken, err := h.authService.SwitchOrg(userID, req.OrgID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"access_token": accessToken})
}

// Logout godoc
// @Summary      Logout user
// @Description  Logs out the user by invalidating the refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        token  body  object{refresh_token=string}  true  "Refresh token"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	err := h.authService.Logout(req.RefreshToken)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Logged out successfully"})
}

// Refresh godoc
// @Summary      Refresh tokens
// @Description  Generates new access and refresh tokens using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        token  body  object{refresh_token=string}  true  "Refresh token"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	accessToken, refreshToken, err := h.authService.Refresh(req.RefreshToken)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
