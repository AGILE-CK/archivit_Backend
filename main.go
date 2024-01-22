package main

import (
	"archivit_Backend/src/domain/auth"
	"archivit_Backend/src/domain/ping"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {

	//dbConfig := db.DataSource{}
	//
	//dataSource := dbConfig.GetDataSource()
	//defer dataSource.Close()
	//db.GormMigrate()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)

	r := gin.Default()

	r.GET("/ping", ping.RequestPing)
	r.GET("/auth/google/login", auth.GoogleLoginHandler)
	r.GET("/auth/google/callback", auth.GoogleAuthCallback)

	r.Run(":8080")
}
