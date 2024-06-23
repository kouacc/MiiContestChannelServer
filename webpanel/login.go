package webpanel

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func (w *WebPanel) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}



func (w *WebPanel) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}



func (w *WebPanel) AdminPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}
