package auth

import (
	"archivit_Backend/src/domain/user"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginHandler godoc
// @title Login
// @Summary login user
// @Schemes
// @Description login user
// @Accept json
// @Produce json
// @Param request body LoginRequest true "request"
// @Success 200 {string} string "User created successfully"
// @Failure 400 {string} string "Invalid user data"
// @Router /auth/login [post]
func LoginHandler(c *gin.Context) {
	//var u entity.User
	var loginRequest LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}

	userInDb, _ := user.FindUserByEmail(loginRequest.Email)
	if userInDb == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "There is no Account with this email address"})
		return
	}

	hash := sha256.Sum256([]byte(loginRequest.Password))
	hashedPassword := hex.EncodeToString(hash[:])

	if loginRequest.Email == userInDb.Email && hashedPassword == userInDb.Password && userInDb.LoginType == "NORMAL" {
		tokenString, err := CreateToken(loginRequest.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}
