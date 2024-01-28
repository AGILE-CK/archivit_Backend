package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

func InitSecretKey() {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		log.Fatalln("SECRET_KEY is required")
	}
	secretKey = []byte(secret)
}

// array of strings
var skipUrls = []string{
	"/swagger/index.html",
	"/swagger-ui/*any",
	"/auth",
	"/auth/google",
}

func CreateToken(email string) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		},
	)

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		//if strings.HasPrefix(c.Request.URL.Path, "/swagger") || strings.HasPrefix(c.Request.URL.Path, "/swagger-ui") {
		//	c.Next()
		//	return
		//}

		for _, skipUrl := range skipUrls {
			if strings.HasPrefix(c.Request.URL.Path, skipUrl) {
				c.Next()
				return
			}
		}
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		tokenString = tokenString[len("Bearer "):]

		err := verifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// If token is valid, proceed with the request
		c.Next()
	}
}
