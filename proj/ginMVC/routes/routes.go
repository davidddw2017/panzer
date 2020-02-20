package routes

import (
	"github.com/davidddw2017/panzer/proj/ginMvc/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/", controller.IndexHome)
	router.GET("/index", controller.IndexHome)
	router.GET("/users/:id", controller.UserGet)
	router.GET("/users", controller.UserGetList)
	router.POST("/users", controller.UserPost)
	router.PUT("/users/:id", controller.UserPut)
	// router.PATCH("/users/:id", controllers.UserPut)
	router.DELETE("/users/:id", controller.UserDelete)
}
