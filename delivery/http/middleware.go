package http

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"gotest/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware menyisipkan/membaca Unique Request ID untuk distributed tracing
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		c.Set(string(utils.RequestIDKey), reqID)
		c.Header("X-Request-ID", reqID)

		c.Next()
	}
}

// StructuredLoggerMiddleware mencatat seluruh request HTTP ke log format JSON (slog)
func StructuredLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		reqID, _ := c.Get(string(utils.RequestIDKey))

		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		// Log terstruktur format JSON
		logAttrs := []slog.Attr{
			slog.String("request_id", reqID.(string)),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.Int("status", statusCode),
			slog.Float64("latency_ms", float64(latency.Microseconds())/1000.0),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
		}

		if statusCode >= 500 {
			slog.LogAttrs(c.Request.Context(), slog.LevelError, "HTTP Server Error", logAttrs...)
		} else if statusCode >= 400 {
			slog.LogAttrs(c.Request.Context(), slog.LevelWarn, "HTTP Client Error", logAttrs...)
		} else {
			slog.LogAttrs(c.Request.Context(), slog.LevelInfo, "HTTP Request Completed", logAttrs...)
		}
	}
}

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

		// Simpan user_id ke context Gin
		userIDFloat, ok := claims["user_id"].(float64)
		if ok {
			c.Set("user_id", uint(userIDFloat))
		}
		c.Set("email", claims["email"])

		c.Next()
	}
}
