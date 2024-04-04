package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strings"
	"time"
	"website-telemetry-demo/api"
	"website-telemetry-demo/cmd/app/middlewares"
	"website-telemetry-demo/configs"
)

func StartApp(config configs.StaticConfig) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router = api.HandleAPI(router)

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
		c.SetCookie("user_auth_token", "", -1, "/", "", true, true) // remove cookie

		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title": "Main website",
		})
	})

	group := router.Group("/", middlewares.RequireAuth())
	group.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	group.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	err := router.Run(net.JoinHostPort(config.Host, config.Port))
	if err != nil {
		return
	}
}
