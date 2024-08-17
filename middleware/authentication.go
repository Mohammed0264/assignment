package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

func AuthMiddleWareMember(role string) gin.HandlerFunc {
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "please login"})
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNoCookie
			}
			return secretKey, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access token expired"})
			return
		}

		claims, ok := token.Claims.(*customClaims)
		if role != "" {
			if claims.Role != "Admin" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "you do not have permission"})
				return
			}
		}

		if ok && token.Valid {
			c.Next()
		} else {

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "login please"})
			return
		}
	}

}
func LogOut(c *gin.Context) {
	var userinfo UserInformation
	err := c.Bind(&userinfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("refreshToken", "", -1, "/", "", false, true)
}
