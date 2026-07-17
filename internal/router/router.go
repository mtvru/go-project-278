package router

import (
	"net/http"
	"reflect"
	"strings"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mtvru/go-project-278/internal/handler"
)

func New(links *handler.LinkHandler, visits *handler.VisitHandler, allowOrigins []string) *gin.Engine {
	registerJSONFieldNames()

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Range"},
		ExposeHeaders:    []string{"Content-Range"},
		AllowCredentials: true,
	}))
	router.TrustedPlatform = gin.PlatformCloudflare

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.GET("/r/:code", visits.Redirect)

	api := router.Group("/api")
	{
		api.GET("/links", links.List)
		api.POST("/links", links.Create)
		api.GET("/links/:id", links.Get)
		api.PUT("/links/:id", links.Update)
		api.DELETE("/links/:id", links.Delete)
		api.GET("/link_visits", visits.List)
	}

	return router
}

func registerJSONFieldNames() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
