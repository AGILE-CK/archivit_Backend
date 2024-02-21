package file

import (
	"cloud.google.com/go/storage"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"io"
	"net/http"
	"os"
	"strings"
)

// DownloadAllFiles godoc
// @title DownloadAllFiles
// @Summary download all files
// @Schemes
// @Description download all files
// @Produce json
// @Success 200 {string} string "All files downloaded successfully"
// @Failure 400 {string} string "Invalid JWT"
// @Router /file/download/all [get]
// @Security ApiKeyAuth
func DownloadAllFiles(c *gin.Context) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	secretKey := os.Getenv("SECRET_KEY")
	bucketName := os.Getenv("BUCKET_NAME")

	jwtToken := c.GetHeader("Authorization")
	jwtToken = jwtToken[len("Bearer "):]

	// JWT에서 email 추출
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JWT: " + err.Error()})
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email := claims["email"].(string)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client: " + err.Error()})
			return
		}

		bucket := client.Bucket(bucketName)

		it := bucket.Objects(ctx, &storage.Query{Prefix: email + "/"})
		for {
			attrs, err := it.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if strings.HasSuffix(attrs.Name, "/") {
				continue
			}

			// Create the directories in the path
			dirs := strings.Split(attrs.Name, "/")
			if len(dirs) > 1 {
				if err := os.MkdirAll(strings.Join(dirs[:len(dirs)-1], "/"), os.ModePerm); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}

			// Create the file
			f, err := os.Create(attrs.Name)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer f.Close()

			// Download the object
			rc, err := bucket.Object(attrs.Name).NewReader(ctx)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer rc.Close()

			if _, err := io.Copy(f, rc); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "All files downloaded successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JWT"})
	}

}
