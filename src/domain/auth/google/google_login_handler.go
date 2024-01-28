package google

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"time"
)

var googleOauthConfig = oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	ClientID:     "584965141712-eku96vnto2vr7t4bk584kkf7q4mer4hn.apps.googleusercontent.com",
	ClientSecret: "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",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
} // todo: change redirect url & osenv

func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour)

	b := make([]byte, 16)
	_, _ = rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
}

// GoogleLoginHandler godoc
// @Summary google login handler
// @Description google login handler
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"ok"
// @Router /auth/google/login [get]
func GoogleLoginHandler(c *gin.Context) {
	state := generateStateOauthCookie(c.Writer)
	url := googleOauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)

}
