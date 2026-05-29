package controllers

import (
	"sport_platform/application/handlers"
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
)

func ClubController(engine *gin.Engine, wrapper *service_wrapper.Wrapper) {
	routerGroup := engine.Group("/clubs")
	routerGroup.POST("/create", func(context *gin.Context) {
		handlers.CreateClubHandler(context, wrapper)
	})
	routerGroup.GET("/:id/workouts", func(context *gin.Context) {
		handlers.GetClubWorkoutsHandler(context, wrapper)
	})
	routerGroup.DELETE("/:id", func(context *gin.Context) {
		handlers.DeleteClubHandler(context, wrapper)
	})
	routerGroup.GET("/", func(context *gin.Context) {
		handlers.GetClubsHandler(context, wrapper)
	})
	routerGroup.GET("/:id", func(context *gin.Context) {
		handlers.GetClubHandler(context, wrapper)
	})
}
