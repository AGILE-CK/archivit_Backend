package record

import (
	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"os"
)

// CreateRecord godoc
// @title CreateRecord
// @Summary create record
// @Schemes
// @Description create record
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "file m4a"
// @Success 200 {string} string "Record created successfully"
// @Failure 400 {string} string "Invalid JWT"
// @Router /record/create [post]
// @Security ApiKeyAuth
func CreateRecord(c *gin.Context) {

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

		ctx := context.Background()
		client, err := storage.NewClient(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client: " + err.Error()})
			return
		}

		bucket := client.Bucket(bucketName)

		emailDir := bucket.Object(email + "/")
		if _, err := emailDir.Attrs(ctx); err != nil {
			wc := emailDir.NewWriter(ctx)
			wc.Close()
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file: " + err.Error()})
			return
		}

		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open file: " + err.Error()})
			return
		}
		defer f.Close()

		wc := bucket.Object(email + "/" + file.Filename).NewWriter(ctx)
		wc.ContentType = "audio/m4a"
		wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
		if _, err := io.Copy(wc, f); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file: " + err.Error()})
			return
		}
		wc.Close()

		c.JSON(http.StatusOK, gin.H{"message": "File created successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JWT"})
	}
}
