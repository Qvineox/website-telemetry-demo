package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
	"website-telemetry-demo/api/middlewares"
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
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		sessionUUID := uuid.New()

		c.SetCookie("user_session_token", sessionUUID.String(), 86400, "/", "", true, true)
		c.SetCookie("user_name", payload.Username, 86400, "/", "", true, true)
	})

	api.POST("/logout", func(c *gin.Context) {
		c.SetCookie("user_session_token", "", -1, "/", "", true, true)

		c.Redirect(http.StatusFound, "/login")
		c.Abort()
	}, middlewares.RequireAuth())

	monitoringGroup := api.Group("/monitoring", middlewares.RequireAuth())
	monitoringGroup.POST("/event", func(c *gin.Context) {
		var payload eventPayload

		err := c.ShouldBindJSON(&payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		payload.Timestamp = time.Now()
		payload.SessionUUID = c.GetString("user_session_token")
		payload.Username = c.GetString("user_name")

		slog.Info(fmt.Sprintf("registered event: %s", payload.String()))
	})

	return router
}

type authPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type eventPayload struct {
	Element   string `json:"element" binding:"required"`
	EventType string `json:"event_type" binding:"required"`
	Message   string `json:"message" binding:"required"`

	// evaluated on the backend
	SessionUUID string    `json:"-"`
	Username    string    `json:"-"`
	Timestamp   time.Time `json:"-"`
}

func (payload *eventPayload) String() string {
	return fmt.Sprintf("session: %s, message: %s", payload.SessionUUID, payload.Message)
}
