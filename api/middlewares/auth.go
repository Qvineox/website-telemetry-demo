package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionTokenCookie, err := c.Cookie("user_session_token")
		if err != nil || len(sessionTokenCookie) == 0 {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()

			return
		}

		sessionUserCookie, err := c.Cookie("user_name")
		if err != nil || len(sessionUserCookie) == 0 {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()

			return
		}

		c.Set("user_session_token", sessionTokenCookie)
		c.Set("user_name", sessionUserCookie)

		c.Next()
	}
}
