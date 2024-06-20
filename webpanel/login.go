package webpanel

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const QueryUser = `SELECT password FROM panel_user WHERE username = $1`

func (w *WebPanel) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (w *WebPanel) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var expectedHash string
	err := w.Pool.QueryRow(w.Ctx, QueryUser, username).Scan(&expectedHash)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"Error": "Database error",
		})
		return
	}

	hashBytes := sha512.Sum512(append(w.Salt, []byte(password)...))
	fmt.Println(hex.EncodeToString(hashBytes[:]))
	if hex.EncodeToString(hashBytes[:]) != expectedHash {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"Error": "Authentication failed",
		})
		return
	}

	// Now that we know the user is valid, generate a JWT.
	claims := &JWTClaims{
		Username: username,
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
	c.Redirect(http.StatusMovedPermanently, "/panel/contests")
}

func (w *WebPanel) AdminPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}
