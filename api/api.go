package api

import (
	"github.com/gin-gonic/gin"
)

func HandleAPI(router *gin.Engine) *gin.Engine {
	api := router.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})

	api.POST("/login", func(c *gin.Context) {
		var payload authPayload

		err := c.ShouldBindJSON(&payload)
		if err != nil {
			c.JSON(401, gin.H{"message": err.Error()})
			return
		} else if len(payload.Username) == 0 || len(payload.Password) == 0 {
			c.JSON(401, gin.H{"message": "username or password are empty"})
			return
		}

		c.SetCookie("user_auth_token", "user_cookie_value", 86400, "/", "", true, true)
	})

	return router
}

type authPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
