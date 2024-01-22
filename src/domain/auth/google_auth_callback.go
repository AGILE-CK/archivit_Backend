package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"io"
	"log"
	"net/http"
)

func getGoogleUserInfo(code string) ([]byte, error) { // 17

	token, err := googleOauthConfig.Exchange(context.Background(), code) // 18
	if err != nil {                                                      // 19
		return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}

	resp, err := http.Get("./auth/google/login" + token.AccessToken) // 20
	if err != nil {                                                  // 21
		return nil, fmt.Errorf("Failed to Get UserInfo %s\n", err.Error())
	}

	return io.ReadAll(resp.Body) // 23
}

func GoogleAuthCallback(c *gin.Context) {
	oauthstate, _ := c.Request.Cookie("oauthstate") // 12

	if c.Request.FormValue("state") != oauthstate.Value { // 13
		log.Printf("invalid google oauth state cookie:%s state:%s\n", oauthstate.Value, c.Request.FormValue("state"))
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	data, err := getGoogleUserInfo(c.Request.FormValue("code")) // 14
	if err != nil {                                             // 15
		log.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	fmt.Fprint(c.Writer, string(data)) // 16
}
