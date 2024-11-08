package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/razortear88/pokedb/controllers"
)

func TypeRoute(route *gin.RouterGroup) {
	typeRoute := route.Group("type")

	typeRoute.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create_type.html", gin.H{
			"title": "Create Type",
		})
	})
	typeRoute.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "list_type.html", gin.H{
			"types": controllers.GetAllTypes(),
			"title": "List Type",
		})
	})
	typeRoute.GET("/:typeName", func(c *gin.Context) {
		c.HTML(http.StatusOK, "detail_type.html", gin.H{
			"type":  controllers.GetType(c),
			"title": c.Param("typeName"),
		})
	})
	typeRoute.POST("", controllers.CreateType())
	typeRoute.GET("/:typeName/update", func(c *gin.Context) {
		c.HTML(http.StatusOK, "update_type.html", gin.H{
			"type":  controllers.GetType(c),
			"title": "Update Type " + c.Param("typeName"),
		})
	})
	typeRoute.POST("/:typeName/update", controllers.EditType())
	typeRoute.POST("/:typeName/delete", controllers.DeleteType())
}
