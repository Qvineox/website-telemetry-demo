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

func HandleAPI(router *gin.Engine, e repo.EventsRepo, l repo.LessonsRepo, u repo.UsersRepo) *gin.Engine {
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

	{
		usersGroup := api.Group("/users")

		usersGroup.GET("/usernames", func(c *gin.Context) {
			users, err := u.GetAllUsernames()
			if err != nil {
				slog.Warn(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

			c.JSON(http.StatusOK, users)
		})

		usersGroup.GET("/sessions", func(c *gin.Context) {
			var payload sessionPayload

			err := c.ShouldBindQuery(&payload)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			clicks, err := u.GetUsernameSessions(payload.Username)
			if err != nil {
				slog.Warn(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

			c.JSON(http.StatusOK, clicks)
		})
	}

	{
		monitoringGroup := api.Group("/monitoring")

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
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

			slog.Info(fmt.Sprintf("registered event: %s", payload.String()))
			c.Status(http.StatusOK)
			return
		})
		monitoringGroup.POST("/mouse-path", func(c *gin.Context) {
			var payload entities.MousePath

			err := c.ShouldBindJSON(&payload)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			payload.Timestamp = time.Now()
			payload.SessionUUID = c.GetString("user_session_token")
			payload.Username = c.GetString("user_name")
			payload.SourceIP = c.ClientIP()

			err = e.SaveMousePath(payload)
			if err != nil {
				slog.Warn(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

			c.Status(http.StatusOK)
			return
		})

		monitoringGroup.GET("/clicks", func(c *gin.Context) {
			var payload heatmapPayload

			err := c.ShouldBindQuery(&payload)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			clicks, err := e.GetClicksByFilter(50000, payload.Location)
			if err != nil {
				slog.Warn(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

			c.JSON(http.StatusOK, clicks)
		})
		monitoringGroup.GET("/mouse-path", func(c *gin.Context) {
			var payload pathPayload

			err := c.ShouldBindQuery(&payload)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			paths, err := e.GetPathByFilter(payload.Session, payload.Location)
			if err != nil {
				slog.Warn(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

			var clicks = make([]entities.ClickData, 0)

			for _, path := range paths {
				for _, pos := range path.Path {
					clicks = append(clicks, entities.ClickData{
						X:     int(pos[0]),
						Y:     int(pos[1]),
						Value: 1,
					})
				}
			}

			c.JSON(http.StatusOK, clicks)
		})
	}

	{
		lessonsGroup := api.Group("/lessons")

		lessonsGroup.POST("/like", func(c *gin.Context) {
			var payload struct {
				ID uint64 `json:"id" binding:"required"`
			}

			err := c.ShouldBindJSON(&payload)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			err = l.LikeLesson(payload.ID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
		})
		lessonsGroup.POST("/dislike", func(c *gin.Context) {
			var payload struct {
				ID uint64 `json:"id" binding:"required"`
			}

			err := c.ShouldBindJSON(&payload)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			err = l.DislikeLesson(payload.ID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
		})
		lessonsGroup.POST("/comment", func(c *gin.Context) {
			var payload struct {
				ID      uint64 `json:"id" binding:"required"`
				Comment string `json:"comment" binding:"required"`
			}

			err := c.ShouldBindJSON(&payload)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			err = l.CommentLesson(payload.ID, entities.Comment{
				Author:    c.GetString("user_name"),
				Text:      payload.Comment,
				LessonID:  &payload.ID,
				CreatedAt: time.Now(),
			})

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
		})
	}

	return router
}

type authPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type heatmapPayload struct {
	Location string `form:"location" binding:"required"`
}

type pathPayload struct {
	Location string `form:"location" binding:"required"`
	Session  string `form:"session" binding:"required"`
}

type sessionPayload struct {
	Username string `form:"username" binding:"required"`
}
