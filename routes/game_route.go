package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/razortear88/pokedb/controllers"
)

func GameRoute(route *gin.RouterGroup) {
	gameRoute := route.Group("game")

	gameRoute.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create_game.html", gin.H{
			"title": "Create Game",
		})
	})
	gameRoute.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "list_game.html", gin.H{
			"games": controllers.GetAllGames(),
		})
	})
	gameRoute.GET("/:gameName", func(c *gin.Context) {
		c.HTML(http.StatusOK, "detail_game.html", gin.H{
			"game":  controllers.GetGame(c),
			"title": c.Param("gameName"),
		})
	})
	gameRoute.POST("", controllers.CreateGame())
	gameRoute.GET("/:gameName/update", func(c *gin.Context) {
		c.HTML(http.StatusOK, "update_game.html", gin.H{
			"game":  controllers.GetGame(c),
			"title": "Update Game " + c.Param("gameName"),
		})
	})
	gameRoute.POST("/:gameName/update", controllers.EditGame())
	gameRoute.POST("/:gameName/delete", controllers.DeleteGame())
}
