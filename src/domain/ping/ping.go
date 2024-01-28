package ping

import (
	"github.com/gin-gonic/gin"
	"os"
)

// RequestPing PingExample godoc
// @title Ping Example API
// @Summary ping example
// @Description do ping hello
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {string} string "pong"
// @Failure 400 {string} string "bad request"
// @Router /ping [get]
func RequestPing(c *gin.Context) {

	temp := os.Getenv("DB_HOST")

	c.JSON(200, gin.H{
		"message": temp,
	})
}
