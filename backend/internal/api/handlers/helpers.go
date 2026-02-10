package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// getUserID safely extracts the userID from the Gin context.
// Returns the user ID and true on success; writes an error response and returns false on failure.
func getUserID(c *gin.Context) (int, bool) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return 0, false
	}
	userID, ok := val.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user session"})
		return 0, false
	}
	return userID, true
}
