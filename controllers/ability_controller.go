package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/razortear88/pokedb/configs"
	"github.com/razortear88/pokedb/models"
	"github.com/razortear88/pokedb/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var abilityCollection *mongo.Collection = configs.GetCollection(configs.DB, "abilities")

// var validate = validator.New()

func CreateAbility() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var ability models.Ability
		defer cancel()
		//validate the request body
		if err := c.ShouldBind(&ability); err != nil {
			c.JSON(http.StatusBadRequest, responses.AbilityResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		if len(c.Request.PostForm["name"]) != 0 {
			ability.Name = c.Request.PostForm["name"][0]
		}
		if len(c.Request.PostForm["description"]) != 0 {
			ability.Description = c.Request.PostForm["description"][0]
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&ability); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.AbilityResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		var abilityExist models.Ability

		err := abilityCollection.FindOne(ctx, bson.M{"name": ability.Name}).Decode(&abilityExist)
		if err == nil {
			c.JSON(http.StatusInternalServerError, responses.AbilityResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Ability Already Exist"}})
			return
		}

		newAbility := models.Ability{
			Name:        ability.Name,
			Description: ability.Description,
		}

		abilityCollection.InsertOne(ctx, newAbility)

		c.Redirect(http.StatusFound, "/ability")
	}
}

func GetAllAbilities() []models.Ability {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var abilities []models.Ability
	defer cancel()
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{"name", 1}})

	results, err := abilityCollection.Find(ctx, filter, opts)

	if err != nil {
		return abilities
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleAbility models.Ability
		if err = results.Decode(&singleAbility); err != nil {
			return abilities
		}
		abilities = append(abilities, singleAbility)
	}

	return abilities
}

func GetAbility(ctx *gin.Context) models.Ability {
	abilityName := ctx.Param("abilityName")
	var ability models.Ability

	abilityCollection.FindOne(ctx, bson.M{"name": abilityName}).Decode(&ability)

	return ability

}

func EditAbility() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		abilityName := c.Param("abilityName")
		var ability models.Ability
		defer cancel()

		//validate the request body
		if err := c.ShouldBind(&ability); err != nil {
			c.JSON(http.StatusBadRequest, responses.AbilityResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if len(c.Request.PostForm["name"]) != 0 {
			ability.Name = c.Request.PostForm["name"][0]
		}

		if len(c.Request.PostForm["description"]) != 0 {
			ability.Description = c.Request.PostForm["description"][0]
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&ability); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.AbilityResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"name": ability.Name, "description": ability.Description}
		_, err := abilityCollection.UpdateOne(ctx, bson.M{"name": abilityName}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.AbilityResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.Redirect(http.StatusFound, "/ability")
	}
}

func DeleteAbility() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		abilityName := c.Param("abilityName")
		defer cancel()

		result, err := abilityCollection.DeleteOne(ctx, bson.M{"name": abilityName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.AbilityResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.AbilityResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Ability with specified Name not found!"}},
			)
			return
		}

		c.Redirect(http.StatusFound, "/ability")
	}
}
