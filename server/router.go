package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()

	router.NoMethod()

	router.NoRoute(func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "")
	})

	router.Use(middlewares...)

	{
		router.GET("/metrics", PrometheusHandler())
	}

	// group
	sampleGroup := router.Group("/sample")
	sampleGroup.Use(
	// unique middleware
	)

	return router
}
