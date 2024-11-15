package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/razortear88/pokedb/controllers"
)

func MoveRoute(route *gin.RouterGroup) {
	moveRoute := route.Group("move")

	moveRoute.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create_move.html", gin.H{
			"title": "Create Move",
			"types": controllers.GetAllTypes(),
		})
	})
	moveRoute.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "list_move.html", gin.H{
			"moves": controllers.GetAllMoves(),
			"title": "List Move",
		})
	})
	moveRoute.GET("/:moveName", func(c *gin.Context) {
		c.HTML(http.StatusOK, "detail_move.html", gin.H{
			"move":  controllers.GetDetailedMove(c),
			"title": c.Param("moveName"),
		})
	})
	moveRoute.POST("", controllers.CreateMove())
	moveRoute.GET("/:moveName/update", func(c *gin.Context) {
		c.HTML(http.StatusOK, "update_move.html", gin.H{
			"move":  controllers.GetMove(c),
			"title": "Update Move " + c.Param("moveName"),
			"types": controllers.GetAllTypes(),
		})
	})
	moveRoute.POST("/:moveName/update", controllers.EditMove())
	moveRoute.POST("/:moveName/delete", controllers.DeleteMove())
}
