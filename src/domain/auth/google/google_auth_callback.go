package google

import (
	"archivit_Backend/src/db/entity"
	"archivit_Backend/src/domain/auth"
	"archivit_Backend/src/domain/user"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"io"
	"log"
	"net/http"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func getGoogleUserInfo(code string) ([]byte, error) {

	token, err := googleOauthConfig.Exchange(context.Background(), code) // 18
	if err != nil {                                                      // 19
		return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}

	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get UserInfo %s\n", err.Error())
	}

	return io.ReadAll(resp.Body)
}

func GoogleAuthCallback(c *gin.Context) {
	oauthstate, _ := c.Request.Cookie("oauthstate")

	if c.Request.FormValue("state") != oauthstate.Value {
		log.Printf("invalid google oauth state cookie:%s state:%s\n", oauthstate.Value, c.Request.FormValue("state"))
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	data, err := getGoogleUserInfo(c.Request.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	fmt.Fprint(c.Writer, string(data))

	//{
	//	"id": "103564838690673927306",
	//	"email": "sungu0804@gmail.com",
	//	"verified_email": true,
	//	"picture": "https://lh3.googleusercontent.com/a-/ALV-UjUSFbSzf50TrMb6EYU-0w5h9cZw3xHWsReUaNXmGOLe=s96-c"
	//}

	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}

	if err := json.Unmarshal(data, &googleUser); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing Google user data"})
		return
	}

	userInDb, _ := user.FindUserByEmail(googleUser.Email)

	// If the user does not exist, create a new user
	if userInDb == nil {
		userInDb = &entity.User{
			Email:     googleUser.Email,
			Password:  "", // Password is empty
			LoginType: "GOOGLE",
		}
		if err := user.SaveUser(userInDb); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user"})
			return
		}
	}

	// Generate a token for the user
	tokenString, err := auth.CreateToken(googleUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
