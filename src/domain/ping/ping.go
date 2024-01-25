package ping

import (
	"github.com/gin-gonic/gin"
	"os"
)

// RequestPing PingExample godoc
// @title Ping Example API
// @Summary ping example
// @Schemes
// @Description do ping hello
// @Accept json
// @Produce json
// @Success 200 {string} string "pong"
// @Router /ping [get]
func RequestPing(c *gin.Context) {

	temp := os.Getenv("DB_HOST")

	c.JSON(200, gin.H{
		"message": temp,
	})
}
