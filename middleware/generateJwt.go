package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

// Define custom claims
type customClaims struct {
	Id   uint
	Role string
	jwt.RegisteredClaims
}

type UserInformation struct {
	UserName string `json:"user-name"`
	Id       uint
	Role     string
}

// CreateToken generates a JWT token for the user
func CreateToken(user UserInformation) (string, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	claims := customClaims{
		Id:   user.Id,
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UserName,
			Audience:  jwt.ClaimStrings{"assignment"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
func GenerateRefreshToken(user UserInformation) (string, error) {
	refreshKey := []byte(os.Getenv("REFRESH_KEY"))
	claims := customClaims{
		Id:   user.Id,
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UserName,
			Audience:  jwt.ClaimStrings{"assignment"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(refreshKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
func RefreshToken(c *gin.Context) {
	refreshKey := []byte(os.Getenv("REFRESH_KEY"))
	var userinfo UserInformation

	cookie, err := c.Cookie("refreshToken")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "token does not exist"})
		return
	}
	token, err := jwt.ParseWithClaims(cookie, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrNoCookie
		}
		return refreshKey, nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	_, ok := token.Claims.(*customClaims)
	if !ok || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid token"})
		return
	}
	userinfo.Id = token.Claims.(*customClaims).Id
	userinfo.Role = token.Claims.(*customClaims).Role
	userinfo.UserName = token.Claims.(*customClaims).Subject
	newAccessToken, err := CreateToken(userinfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": newAccessToken})
}

//return signedToken, nil
