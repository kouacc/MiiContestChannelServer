package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			// We can't redirect off an Unauthorized status code.
			c.Redirect(http.StatusPermanentRedirect, "/panel/login")
			c.Abort()
			return
		}

		claims, err := VerifyToken(tokenString)
		if err != nil {
			c.Redirect(http.StatusPermanentRedirect, "/panel/login")
			c.Abort()
			return
		}

		if username, ok := claims["username"].(string); ok {
			c.Set("username", username)
		} else {
			c.Redirect(http.StatusPermanentRedirect, "/panel/login")
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}
