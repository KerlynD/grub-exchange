package handlers

import (
	"grub-exchange/internal/models"
	"grub-exchange/internal/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	user, token, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setAuthCookie(c, token)
	c.JSON(http.StatusCreated, gin.H{"user": user, "token": token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	user, token, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	setAuthCookie(c, token)
	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "grub_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, _ := c.Get("userID")

	user, err := h.authService.GetMe(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func setAuthCookie(c *gin.Context, token string) {
	isProduction := os.Getenv("FRONTEND_URL") != "" && os.Getenv("FRONTEND_URL") != "http://localhost:3000"

	cookie := &http.Cookie{
		Name:     "grub_token",
		Value:    token,
		MaxAge:   7 * 24 * 60 * 60, // 7 days
		Path:     "/",
		HttpOnly: true,
	}

	if isProduction {
		// Cross-domain: Vercel frontend → Fly.io backend
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	} else {
		// Local dev: same origin
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	}

	http.SetCookie(c.Writer, cookie)
}
