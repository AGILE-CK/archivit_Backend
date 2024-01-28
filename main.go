package main

import (
	"archivit_Backend/docs"
	"archivit_Backend/src/db"
	"archivit_Backend/src/domain/auth"
	"archivit_Backend/src/domain/auth/google"
	"archivit_Backend/src/domain/ping"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
)

func setupSwagger(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger/index.html")
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// @title ARCHIVIT API
func main() {

	dbConfig := db.DataSource{}

	dataSource := dbConfig.MakeDataSource()
	defer dataSource.Close()
	db.GormMigrate()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)

	r := gin.Default()

	docs.SwaggerInfo.BasePath = ""
	setupSwagger(r)

	r.GET("/ping", ping.RequestPing)
	r.GET("/auth/google/login", google.GoogleLoginHandler)
	r.GET("/auth/google/callback", google.GoogleAuthCallback)

	r.POST("/auth/signup", auth.RegisterHandler)
	r.POST("/auth/login", auth.LoginHandler)

	r.Run(":" + port)
}
