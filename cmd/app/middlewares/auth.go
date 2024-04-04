package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("user_auth_token")
		if err != nil || len(cookie) == 0 {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()

			return
		}

		c.Next()
	}
}
