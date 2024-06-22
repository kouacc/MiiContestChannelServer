package webpanel

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	oauth2State  string
	oauth2Config *oauth2.Config
)



func NewAppAuthConfig(ctx context.Context) (*AppAuthConfig, error) {
    provider, err := oidc.NewProvider(ctx, "https://sso.riiconnect24.net/application/o/cmoc-localhost-testing/.well-known/openid-configuration")
    if err != nil {
        return nil, err
    }

    oauth2Config := &oauth2.Config{
        ClientID:     "pddfSBVj9B3GkZZgo2GFfaGsxU6YUP6rZex32XI2",
        ClientSecret: "WReS7ZCAbkx0MyoaWd9APVSpjYYFgWUYCLPvuDtF3XoQoju1gpVuQbeEk6MAMorQPMFP2OEFM8TqXfvzhP7vNMQ2ACou9A8r5d1dcKw2A8axmXRRWunMG80d6k32Wt37",
        RedirectURL:  "http://localhost:9011/panel/authorize",
        Scopes:       []string{"openid", "profile"},
        Endpoint:     provider.Endpoint(),
    }

    return &AppAuthConfig{
        OAuth2Config: oauth2Config,
        Provider:     provider,
    }, nil
}

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
	
	http.Redirect(c.Writer, c.Request, AuthConfig.OAuth2Config.AuthCodeURL(state), http.StatusFound)
}



func (w *WebPanel) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}



func (w *WebPanel) AdminPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}
