package webpanel

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
	"github.com/golang-jwt/jwt/v5"
)

var (
	oauth2State  string
	oauth2Config *oauth2.Config
)

func init() {
	oauth2Config = &oauth2.Config{
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://your-authentik-instance/authorize",
			TokenURL: "https://your-authentik-instance/token",
		},
	}
	oauth2State = "help-me"
}
	

func (w *WebPanel) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (w *WebPanel) Login(c *gin.Context) {
	url := oauth2Config.AuthCodeURL(oauth2State)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (w *WebPanel) AuthCallback(c *gin.Context) {
	state := c.Query("state")
	if state != oauth2State {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"Error": "Invalid state parameter",
		})
		return
	}

	code := c.Query("code")
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"Error": "Failed to exchange token",
		})
		return
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"Error": "No id_token field in oauth2 token",
		})
		return
	}

	// Parse the ID Token to get user information
	claims := &JWTClaims{}
	_, err = jwt.ParseWithClaims(idToken, claims, func(token *jwt.Token) (interface{}, error) {
		// Use the public key of your Authentik instance to verify the token
		return []byte("your-public-key"), nil
	})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"Error": "Failed to parse ID token",
		})
		return
	}

	// Use the claims to get user information and create a session
	c.SetCookie("token", idToken, 3600, "", "", false, true)
	c.Redirect(http.StatusMovedPermanently, "/panel/contests")
}


func (w *WebPanel) AdminPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}
