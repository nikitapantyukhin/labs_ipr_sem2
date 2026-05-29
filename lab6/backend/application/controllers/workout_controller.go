package controllers

import (
	"sport_platform/application/handlers"
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
)

func WorkoutController(engine *gin.Engine, wrapper *service_wrapper.Wrapper) {
	workoutsGroup := engine.Group("/workouts")
	workoutsGroup.DELETE("/:workout_id", func(context *gin.Context) {
		handlers.DeleteWorkoutHandler(context, wrapper)
	})
	workoutsGroup.POST("/create", func(context *gin.Context) {
		handlers.CreateWorkoutHandler(context, wrapper)
	})
	workoutsGroup.PUT("/update", func(context *gin.Context) {
		handlers.UpdateWorkoutHandler(context, wrapper)
	})
}
