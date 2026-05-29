package controllers

import (
	"sport_platform/application/handlers"
	"sport_platform/internal/service_wrapper"

	"github.com/gin-gonic/gin"
)

func UserController(engine *gin.Engine, wrapper *service_wrapper.Wrapper) {
	routerGroup := engine.Group("/users")
	routerGroup.POST("/create", func(context *gin.Context) {
		handlers.CreateUserHandler(context, wrapper)
	})
	routerGroup.POST("/login", func(context *gin.Context) {
		handlers.LoginHandler(context, wrapper)
	})
	routerGroup.PUT("/update", func(context *gin.Context) {
		handlers.UpdateUserHandler(context, wrapper)
	})
	routerGroup.GET("/", func(context *gin.Context) {
		handlers.GetUserHandler(context, wrapper)
	})
	routerGroup.DELETE("/delete", func(context *gin.Context) {
		handlers.DeleteUserHandler(context, wrapper)
	})
}
