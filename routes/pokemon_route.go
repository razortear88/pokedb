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

	pokemonRoute.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "list_pokemon.html", gin.H{
			"pokemons": controllers.GetAllPokemons(),
			"title":    "List Pokemon",
		})
	})
	pokemonRoute.GET("/:pokemonName", func(c *gin.Context) {
		c.HTML(http.StatusOK, "detail_pokemon.html", gin.H{
			"pokemon": controllers.GetDetailedPokemon(c),
		})
	})

	pokemonRoute.GET("/:pokemonName/update", func(c *gin.Context) {
		pokemon := controllers.GetPokemon(c)
		if pokemon.Name == "" {
			c.HTML(http.StatusOK, "update_pokemon.html", gin.H{
				"pokemon": pokemon,
			})
		}
		var pokemonType2 string

		if len(pokemon.Type) == 2 {
			pokemonType2 = pokemon.Type[1]
		} else {
			pokemonType2 = ""
		}

		c.HTML(http.StatusOK, "update_pokemon.html", gin.H{
			"pokemon":      pokemon,
			"pokemonType1": pokemon.Type[0],
			"pokemonType2": pokemonType2,
			"types":        controllers.GetAllTypes(),
		})
	})

	pokemonRoute.POST("", controllers.CreatePokemon())
	pokemonRoute.POST("/:pokemonName/update", controllers.EditPokemon())
	pokemonRoute.POST("/:pokemonName/delete", controllers.DeletePokemon())
	pokemonRoute.GET("/api/list", controllers.GetApiPokemons)
}
