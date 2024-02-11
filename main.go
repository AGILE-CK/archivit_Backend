package main

import (
	"archivit_Backend/docs"
	"archivit_Backend/src/db"
	"archivit_Backend/src/domain/auth"
	"archivit_Backend/src/domain/auth/google"
	"archivit_Backend/src/domain/ping"
	"cloud.google.com/go/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/net/context"
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
// @version latest
// @Content-Type application/json
// @description This is server for Archivit API.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	google.InitConfig()
	auth.InitSecretKey()

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

	router := gin.Default()
	router.Use(cors.Default())

	docs.SwaggerInfo.BasePath = ""
	setupSwagger(router)

	router.Use(auth.TokenAuthMiddleware())

	router.GET("/ping", ping.RequestPing)
	router.GET("/auth/google/login", google.GoogleLoginHandler)
	router.GET("/auth/google/callback", google.GoogleAuthCallback)

	router.POST("/auth/signup", auth.RegisterHandler)
	router.POST("/auth/login", auth.LoginHandler)

	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	//projectID := "primeval-span-410215"

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name for the new bucket.
	bucketName := "primeval-span-410215.appspot.com"

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	//// Creates the new bucket.
	//
	//if err := bucket.Create(ctx, projectID, nil); err != nil {
	//	log.Fatalf("Failed to create bucket: %v", err)
	//}

	// write a new file in bucket
	wc := bucket.Object("testfile").NewWriter(ctx)
	wc.ContentType = "text/plain"
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	wc.Write([]byte("hello world\n"))
	wc.Close()

	router.Run(":" + port)
}
