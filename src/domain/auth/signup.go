package auth

import (
	"archivit_Backend/src/db/entity"
	"archivit_Backend/src/domain/user"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func checkEmailForm(email string) bool {
	return strings.Contains(email, "@")
}

func checkPasswordForm(password string) bool {
	return len(password) >= 8
}

func RegisterHandler(c *gin.Context) {
	var u entity.User

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}

	userInDb, _ := user.FindUserByEmail(u.Email)

	if u.Email == userInDb.Email {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	if checkEmailForm(u.Email) && checkPasswordForm(u.Password) {

		hash := sha256.Sum256([]byte(u.Password))
		u.Password = hex.EncodeToString(hash[:])
		u.LoginType = "NORMAL"

		err := user.SaveUser(&u)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "User created successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
	}
}
