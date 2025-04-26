package utils

import (
	"github.com/gin-gonic/gin"
)

// ContextUserKey is the key used to store the user ID in the context
const ContextUserKey = "userID"

// GetUserIDFromContext retrieves the user ID from the Gin context
// This function assumes the user ID was set by the JWT middleware
func GetUserIDFromContext(c *gin.Context) uint {
	// Get the user ID from the context
	value, exists := c.Get(ContextUserKey)
	if !exists {
		return 0
	}

	// Convert the value to uint
	userID, ok := value.(uint)
	if !ok {
		return 0
	}

	return userID
}
