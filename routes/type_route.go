package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/razortear88/pokedb/controllers"
)

func TypeRoute(router *gin.Engine) {
	router.GET("/types", controllers.GetAllTypes())
	router.GET("/type/:typeName", controllers.GetAType())
	router.POST("/type", controllers.CreateType())
	router.PUT("/type/:typeName", controllers.EditAType())
	router.DELETE("/type/:typeName", controllers.DeleteAType())
}
