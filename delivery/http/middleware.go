package http

import (
	"net/http"
	"strings"

	"gotest/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware memproteksi route dengan memverifikasi Bearer Token JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Token otorisasi diperlukan (Bearer <token>)",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		// Simpan user_id ke context Gin agar bisa dipakai di handler berikutnya
		userIDFloat, ok := claims["user_id"].(float64)
		if ok {
			c.Set("user_id", uint(userIDFloat))
		}
		c.Set("email", claims["email"])

		c.Next()
	}
}
