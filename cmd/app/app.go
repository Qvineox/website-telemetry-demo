package app

import (
	"github.com/a-h/templ"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net"
	"strings"
	"time"
	"website-telemetry-demo/configs"
	"website-telemetry-demo/web/templates"
)

func StartApp(config configs.StaticConfig) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

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

	router.GET("/", func(context *gin.Context) {
		templ.Handler(templates.Index(time.Now())).ServeHTTP(context.Writer, context.Request)
	})

	err := router.Run(net.JoinHostPort(config.Host, config.Port))
	if err != nil {
		return
	}
}
