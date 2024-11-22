package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/razortear88/pokedb/models"
	"github.com/razortear88/pokedb/responses"
)

// var pokemonCollection *mongo.Collection = configs.GetCollection(configs.DB, "pokemons")

func CreatePokemon() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pokemon models.Pokemon
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		//validate the request body
		if err := c.ShouldBind(&pokemon); err != nil {
			c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

	}
}
