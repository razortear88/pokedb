package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/razortear88/pokedb/controllers"
)

func AbilityRoute(route *gin.RouterGroup) {
	abilityRoute := route.Group("ability")

	abilityRoute.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create_ability.html", gin.H{
			"title": "Create Ability",
		})
	})
	abilityRoute.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "list_ability.html", gin.H{
			"abilities": controllers.GetAllAbilities(),
			"title":     "List Ability",
		})
	})
	abilityRoute.GET("/:abilityName", func(c *gin.Context) {
		c.HTML(http.StatusOK, "detail_ability.html", gin.H{
			"ability": controllers.GetAbility(c),
			"title":   c.Param("abilityName"),
		})
	})
	abilityRoute.POST("", controllers.CreateAbility())
	abilityRoute.GET("/:abilityName/update", func(c *gin.Context) {
		c.HTML(http.StatusOK, "update_ability.html", gin.H{
			"ability": controllers.GetAbility(c),
			"title":   "Update Ability " + c.Param("abilityName"),
		})
	})
	abilityRoute.POST("/:abilityName/update", controllers.EditAbility())
	abilityRoute.POST("/:abilityName/delete", controllers.DeleteAbility())
}
