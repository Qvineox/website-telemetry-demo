package app

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
	"website-telemetry-demo/api"
	"website-telemetry-demo/api/middlewares"
	"website-telemetry-demo/cmd/app/entities"
	"website-telemetry-demo/cmd/app/repo"
	"website-telemetry-demo/configs"
)

func StartApp(config configs.StaticConfig) {
	gin.SetMode(gin.DebugMode)
	router := gin.New()

	db := prepareDatabaseConnection(config.Database.Host, config.Database.User, config.Database.Password, config.Database.Name, config.Database.Timezone)

	eRepo := repo.NewEventsRepo(db)
	lRepo := repo.NewLessonsRepo(db)

	router = api.HandleAPI(router, eRepo, lRepo)

	router.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(config.AllowedOrigins, ","),
		AllowMethods: []string{"OPTIONS", "GET", "PUT", "PATCH", "DELETE", "POST"},
		AllowHeaders: []string{
			"Accept",
			"Cache-Control",
			"Content-Type",
			"Content-Length",
			"X-CSRF-Token",
			"X-API-Key",
			"Accept-Encoding",
			"Accept-Language",
			"Authorization",
			"X-Forwarded-*",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWildcard:    true,
		MaxAge:           12 * time.Hour,
	}))

	router.Static("/icons", "web/icons")
	router.Static("/styles", "web/styles")
	router.Static("/scripts", "web/scripts")

	router.LoadHTMLGlob("web/templates/*")

	router.GET("/login", func(c *gin.Context) {
		c.SetCookie("user_session_token", "", -1, "/", "", true, true) // remove cookie

		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title": "Main website",
		})
	})

	group := router.Group("/", middlewares.RequireAuth())
	group.GET("/", func(c *gin.Context) {
		events, err := eRepo.GetRecentEvents(10)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		lessons, err := lRepo.GetRecentLessons(5)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"lessons": lessons,
			"events":  events,
		})
	})

	group.GET("/materials", func(c *gin.Context) {
		lessons, err := lRepo.GetAllLessons()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.HTML(http.StatusOK, "lessons.tmpl", gin.H{
			"lessons": lessons,
		})
	})

	group.GET("/materials/:id", func(c *gin.Context) {
		param := c.Param("id")
		id, err := strconv.ParseUint(param, 10, 64)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		lesson, err := lRepo.GetLessonByID(id)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.HTML(http.StatusOK, "lesson.tmpl", lesson)
	})

	group.GET("/profiles", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profiles.tmpl", gin.H{})
	})

	host := net.JoinHostPort(config.Host, config.Port)

	slog.Info(fmt.Sprintf("staring server on: http://%s", host))

	err := router.Run(host)
	if err != nil {
		return
	}
}

func prepareDatabaseConnection(host, user, password, name, timezone string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=%s", host, user, password, name, timezone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&entities.Comment{}, &entities.Lesson{}, &entities.Event{})
	if err != nil {
		panic("failed to migrate schema")
	}

	return db
}
