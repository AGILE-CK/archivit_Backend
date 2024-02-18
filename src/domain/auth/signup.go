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

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterHandler godoc
// @title Register
// @Summary register user
// @Schemes
// @Description register user
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "request"
// @Success 200 {string} string "User created successfully"
// @Failure 400 {string} string "Invalid user data"
// @Router /auth/signup [post]
func RegisterHandler(c *gin.Context) {
	var u entity.User
	var registerRequest RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}

	userInDb, _ := user.FindUserByEmail(registerRequest.Email)

	if userInDb != nil && registerRequest.Email == userInDb.Email {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	if checkEmailForm(registerRequest.Email) && checkPasswordForm(registerRequest.Password) {

		hash := sha256.Sum256([]byte(registerRequest.Password))
		registerRequest.Password = hex.EncodeToString(hash[:])
		u.Email = registerRequest.Email
		u.Password = registerRequest.Password
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
