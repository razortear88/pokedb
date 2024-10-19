package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/razortear88/pokedb/controllers"
	"net/http"
)

func TypeRoute(router *gin.Engine) {
	router.GET("/type/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create_type.html", gin.H{
			"title": "Create Type",
		})
	})
	router.GET("/type", func(c *gin.Context) {
		c.HTML(http.StatusOK, "list_type.html", gin.H{
			"types": controllers.GetAllTypes(),
		})
	})
	router.GET("/type/:typeName", func(c *gin.Context) {
		c.HTML(http.StatusOK, "detail_type.html", gin.H{
			"type": controllers.GetAType(c),
		})
	})
	router.POST("/type", controllers.CreateType())
	router.PUT("/type/:typeName", controllers.EditAType())
	router.GET("/type/:typeName/update", func(c *gin.Context) {
		c.HTML(http.StatusOK, "update_type.html", gin.H{
			"type": controllers.GetAType(c),
		})
	})
	router.POST("/type/:typeName/delete", controllers.DeleteAType())
}
