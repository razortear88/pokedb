package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/razortear88/pokedb/controllers"
)

func PokemonRoute(route *gin.RouterGroup) {
	pokemonRoute := route.Group("pokemon")
	pokemonRoute.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create_pokemon.html", gin.H{
			"title": "Create Pokemon",
			"types": controllers.GetAllTypes(),
		})
	})
	pokemonRoute.POST("", controllers.CreatePokemon())
}
