package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/razortear88/pokedb/configs"
	"github.com/razortear88/pokedb/models"
	"github.com/razortear88/pokedb/responses"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

var typeCollection *mongo.Collection = configs.GetCollection(configs.DB, "types")
var validate = validator.New()

func CreateType() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var typee models.Type
		defer cancel()
		//validate the request body
		if err := c.ShouldBind(&typee); err != nil {
			c.JSON(http.StatusBadRequest, responses.TypeResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		if len(c.Request.PostForm["name"]) != 0 {
			typee.Name = strings.ToUpper(c.Request.PostForm["name"][0])
		}
		if len(c.Request.PostForm["color"]) != 0 {
			typee.Color = c.Request.PostForm["color"][0]
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&typee); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.TypeResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		var typeExist models.Type

		err := typeCollection.FindOne(ctx, bson.M{"name": typee.Name}).Decode(&typeExist)
		if err == nil {
			c.JSON(http.StatusInternalServerError, responses.TypeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Type Already Exist"}})
			return
		}

		newType := models.Type{
			Name:  typee.Name,
			Color: typee.Color,
		}

		result, err := typeCollection.InsertOne(ctx, newType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.TypeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.TypeResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetAllTypes() []models.Type {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var types []models.Type
	defer cancel()
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{"name", 1}})

	results, err := typeCollection.Find(ctx, filter, opts)

	if err != nil {
		return types
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleType models.Type
		if err = results.Decode(&singleType); err != nil {
			return types
		}
		types = append(types, singleType)
	}

	return types
}

func GetAType(ctx *gin.Context) models.Type {
	typeName := ctx.Param("typeName")
	var typee models.Type

	err := typeCollection.FindOne(ctx, bson.M{"name": typeName}).Decode(&typee)
	if err != nil {
	}
	return typee

}

func EditAType() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		typeName := c.Param("typeName")
		var typee models.Type
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&typee); err != nil {
			c.JSON(http.StatusBadRequest, responses.TypeResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&typee); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.TypeResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"name": typee.Name, "color": typee.Color}
		result, err := typeCollection.UpdateOne(ctx, bson.M{"name": typeName}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.TypeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedType models.Type
		if result.MatchedCount == 1 {
			err := typeCollection.FindOne(ctx, bson.M{"name": typeName}).Decode(&updatedType)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.TypeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.TypeResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedType}})
	}
}

func DeleteAType() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		typeName := c.Param("typeName")
		defer cancel()

		result, err := typeCollection.DeleteOne(ctx, bson.M{"name": typeName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.TypeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.TypeResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Type with specified Name not found!"}},
			)
			return
		}

		c.Redirect(http.StatusFound, "/type")
	}
}
