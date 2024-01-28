package auth

import (
	"archivit_Backend/src/db/entity"
	"archivit_Backend/src/domain/user"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginHandler(c *gin.Context) {
	var u entity.User

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}

	userInDb, _ := user.FindUserByEmail(u.Email)

	hash := sha256.Sum256([]byte(u.Password))
	hashedPassword := hex.EncodeToString(hash[:])

	if u.Email == userInDb.Email && hashedPassword == userInDb.Password && userInDb.LoginType == "NORMAL" {
		tokenString, err := CreateToken(u.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}
