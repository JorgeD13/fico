package api

import (
	"net/http"
	"os"
	"strings"

	"fico/gol/db"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func JWTMiddleware(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "dev-secret"
		}
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenMalformed
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// Check revoked
		revoked, _ := db.IsTokenRevoked(gdb, tokenStr)
		if revoked {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// Attach claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userClaims", claims)
		}
		c.Next()
	}
}
