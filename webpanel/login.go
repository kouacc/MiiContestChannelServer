package webpanel

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)



func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func setCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}
	
func (w *WebPanel) StartPanelHandler(c *gin.Context) {
	state, err := randString(16)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}
	
	setCallbackCookie(c.Writer, c.Request, "state", state)
	
	http.Redirect(c.Writer, c.Request, w.AuthConfig.OAuth2Config.AuthCodeURL(state), http.StatusFound)
}

func (w *WebPanel) FinishPanelHandler(c *gin.Context) {
	//TODO: implement response handler for OIDC
	state, err := c.Request.Cookie("state")
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Error": "State cookie not found",
		})
		return
	}

	if c.Request.URL.Query().Get("state") != state.Value {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Error": "State did not match",
		})
		return
	}

	oauth2Token, err := w.AuthConfig.OAuth2Config.Exchange(c, c.Request.URL.Query().Get("code"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	userInfo, err := w.AuthConfig.Provider.UserInfo(c, oauth2.StaticTokenSource(oauth2Token))
		if err != nil {
			http.Error(c.Writer, "Failed to get userinfo: "+err.Error(), http.StatusInternalServerError)
			return
		}

	/* resp := struct {
			OAuth2Token *oauth2.Token
		}{oauth2Token}
		data, err := json.MarshalIndent(resp, "", "    ")
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
		c.Writer.Write(data) */

	//Now that we verified the token, create a JWT token to use with the middleware
	claims := &JWTClaims{
		Username: userInfo.Profile,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("help me"))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"Error": "Failed to create JWT",
		})
		return
	}

	c.SetCookie("token", token, 3600, "", "", false, true)

	//redirect to admin page
	c.Redirect(http.StatusMovedPermanently, "/panel/contests")

}



func (w *WebPanel) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}



func (w *WebPanel) AdminPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}
