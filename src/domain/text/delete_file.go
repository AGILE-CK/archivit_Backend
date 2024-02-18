package text

import (
	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/net/context"
	"net/http"
	"os"
)

// DeleteFile godoc
// @title DeleteFile
// @Summary delete file
// @Schemes
// @Description delete file
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "file txt and m4a"
// @Success 200 {string} string "File deleted successfully"
// @Failure 400 {string} string "Invalid JWT"
// @Router /file/delete [delete]
// @Security ApiKeyAuth
func DeleteFile(c *gin.Context) {
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

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file: " + err.Error()})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		email := claims["email"].(string)

		ctx := context.Background()
		client, err := storage.NewClient(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client: " + err.Error()})
			return
		}

		bucket := client.Bucket(bucketName)

		// Delete the file
		o := bucket.Object(email + "/" + file.Filename)
		if err := o.Delete(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
	}
}
