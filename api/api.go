package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
	"website-telemetry-demo/api/middlewares"
	"website-telemetry-demo/cmd/app/entities"
	"website-telemetry-demo/cmd/app/repo"
)

func HandleAPI(router *gin.Engine, e repo.EventsRepo) *gin.Engine {
	insecure := router.Group("/api")

	insecure.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})

	insecure.POST("/login", func(c *gin.Context) {
		var payload authPayload

		err := c.ShouldBindJSON(&payload)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		sessionUUID := uuid.New()

		c.SetCookie("user_session_token", sessionUUID.String(), 86400, "/", "", true, true)
		c.SetCookie("user_name", payload.Username, 86400, "/", "", true, true)

		_ = e.SaveEvent(entities.Event{
			Element:     "",
			EventType:   "login",
			Message:     "user logged in",
			SessionUUID: sessionUUID.String(),
			Username:    payload.Username,
			SourceIP:    c.RemoteIP(),
			Timestamp:   time.Now(),
		})
	})

	api := insecure.Group("")
	api.Use(middlewares.RequireAuth())

	api.POST("/logout", func(c *gin.Context) {
		_ = e.SaveEvent(entities.Event{
			Element:     "",
			EventType:   "logout",
			Message:     "user logged out",
			SessionUUID: c.GetString("user_session_token"),
			Username:    c.GetString("user_name"),
			SourceIP:    c.RemoteIP(),
			Timestamp:   time.Now(),
		})

		c.SetCookie("user_session_token", "", -1, "/", "", true, true)

		c.Redirect(http.StatusFound, "/login")
		c.Abort()
	})

	monitoringGroup := api.Group("/monitoring", middlewares.RequireAuth())
	monitoringGroup.POST("/event", func(c *gin.Context) {
		var payload entities.Event

		err := c.ShouldBindJSON(&payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		payload.Timestamp = time.Now()
		payload.SessionUUID = c.GetString("user_session_token")
		payload.Username = c.GetString("user_name")
		payload.SourceIP = c.ClientIP()

		err = e.SaveEvent(payload)
		if err != nil {
			slog.Warn(err.Error())
		}

		slog.Info(fmt.Sprintf("registered event: %s", payload.String()))
	})

	return router
}

type authPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
