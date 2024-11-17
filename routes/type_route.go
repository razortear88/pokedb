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
	typeRoute.GET("/:typename", func(c *gin.Context) {
		c.HTML(http.StatusOK, "detail_type.html", gin.H{
			"type":  controllers.GetType(c),
			"title": c.Param("typename"),
		})
	})
	typeRoute.POST("", controllers.CreateType())
	typeRoute.GET("/:typename/update", func(c *gin.Context) {
		c.HTML(http.StatusOK, "update_type.html", gin.H{
			"type":  controllers.GetType(c),
			"title": "Update Type " + c.Param("typename"),
		})
	})
	typeRoute.POST("/:typename/update", controllers.EditType())
	typeRoute.POST("/:typename/delete", controllers.DeleteType())
}
