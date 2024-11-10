package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/razortear88/pokedb/configs"
	"github.com/razortear88/pokedb/models"
	"github.com/razortear88/pokedb/responses"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

		typeCollection.InsertOne(ctx, newType)

		c.Redirect(http.StatusFound, "/type")
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

func GetType(ctx *gin.Context) models.Type {
	typeName := ctx.Param("typeName")
	var typee models.Type

	typeCollection.FindOne(ctx, bson.M{"name": typeName}).Decode(&typee)

	if len(typee.Color) == 4 {
		typee.Color = "#" + strings.Repeat(string(typee.Color[1]), 2) + strings.Repeat(string(typee.Color[2]), 2) + strings.Repeat(string(typee.Color[3]), 2)
	}
	return typee

}

func EditType() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		typeName := c.Param("typeName")
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

		update := bson.M{"name": typee.Name, "color": typee.Color}
		_, err := typeCollection.UpdateOne(ctx, bson.M{"name": typeName}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.TypeResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.Redirect(http.StatusFound, "/type")
	}
}

func DeleteType() gin.HandlerFunc {
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

func GetColor(typeName string) string {
	var typee models.Type

	typeCollection.FindOne(nil, bson.M{"name": typeName}).Decode(&typee)

	return typee.Color
}
