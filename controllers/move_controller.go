package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/razortear88/pokedb/configs"
	"github.com/razortear88/pokedb/models"
	"github.com/razortear88/pokedb/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var moveCollection *mongo.Collection = configs.GetCollection(configs.DB, "moves")

// var validate = validator.New()

func CreateMove() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var move models.Move
		defer cancel()
		//validate the request body
		if err := c.ShouldBind(&move); err != nil {
			c.JSON(http.StatusBadRequest, responses.MoveResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		move.Name = c.Request.PostForm["name"][0]
		move.Category = c.Request.PostForm["category"][0]
		move.TypeName = c.Request.PostForm["typename"][0]

		power, err := strconv.Atoi(c.Request.PostForm["power"][0])
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.MoveResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		move.Power = power

		accuracy, err := strconv.Atoi(c.Request.PostForm["accuracy"][0])
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.MoveResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		move.Accuracy = accuracy

		pp, err := strconv.Atoi(c.Request.PostForm["pp"][0])
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.MoveResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		move.PP = pp
		makecontact, err := strconv.ParseBool(c.Request.PostForm["makecontact"][0])
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.MoveResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		move.MakeContact = makecontact

		move.Effect = c.Request.PostForm["effect"][0]

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&move); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.MoveResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		var moveExist models.Move

		dbErr := moveCollection.FindOne(ctx, bson.M{"name": move.Name}).Decode(&moveExist)
		if dbErr == nil {
			c.JSON(http.StatusInternalServerError, responses.MoveResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Move Already Exist"}})
			return
		}

		newMove := models.Move{
			Name:        move.Name,
			Category:    move.Category,
			TypeName:    move.TypeName,
			Power:       move.Power,
			Accuracy:    move.Accuracy,
			PP:          move.PP,
			MakeContact: move.MakeContact,
			Effect:      move.Effect,
		}

		moveCollection.InsertOne(ctx, newMove)

		c.Redirect(http.StatusFound, "/move")
	}
}

func GetAllMoves() []models.Move {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var moves []models.Move
	defer cancel()
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{"name", 1}})

	results, err := moveCollection.Find(ctx, filter, opts)

	if err != nil {
		return moves
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleMove models.Move
		if err = results.Decode(&singleMove); err != nil {
			return moves
		}
		moves = append(moves, singleMove)
	}

	return moves
}

func GetMove(ctx *gin.Context) models.Move {
	moveName := ctx.Param("moveName")
	var move models.Move

	moveCollection.FindOne(ctx, bson.M{"name": moveName}).Decode(&move)

	return move

}

func EditMove() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		moveName := c.Param("moveName")
		var move models.Move
		defer cancel()

		//validate the request body
		if err := c.ShouldBind(&move); err != nil {
			c.JSON(http.StatusBadRequest, responses.MoveResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		move.Name = c.Request.PostForm["name"][0]
		move.Category = c.Request.PostForm["category"][0]
		move.TypeName = c.Request.PostForm["typename"][0]

		power, err := strconv.Atoi(c.Request.PostForm["power"][0])
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.MoveResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		move.Power = power

		accuracy, err := strconv.Atoi(c.Request.PostForm["accuracy"][0])
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.MoveResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		move.Accuracy = accuracy

		pp, err := strconv.Atoi(c.Request.PostForm["pp"][0])
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.MoveResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		move.PP = pp

		makecontact, err := strconv.ParseBool(c.Request.PostForm["makecontact"][0])
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.MoveResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		move.MakeContact = makecontact

		move.Effect = c.Request.PostForm["effect"][0]

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&move); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.MoveResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{
			"name":        move.Name,
			"category":    move.Category,
			"typename":    move.TypeName,
			"power":       move.Power,
			"accuracy":    move.Accuracy,
			"pp":          move.PP,
			"makecontact": move.MakeContact,
			"effect":      move.Effect}

		_, updateErr := moveCollection.UpdateOne(ctx, bson.M{"name": moveName}, bson.M{"$set": update})
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, responses.MoveResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.Redirect(http.StatusFound, "/move")
	}
}

func DeleteMove() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		moveName := c.Param("moveName")
		defer cancel()

		result, err := moveCollection.DeleteOne(ctx, bson.M{"name": moveName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.MoveResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.MoveResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Move with specified Name not found!"}},
			)
			return
		}

		c.Redirect(http.StatusFound, "/move")
	}
}

func GetDetailedMove(ctx *gin.Context) bson.M {
	moveName := ctx.Param("moveName")

	pipeline := mongo.Pipeline{
		{{"$lookup", bson.D{{"from", "types"}, {"localField", "typename"}, {"foreignField", "name"}, {"as", "lookup"}}}},
		{{"$unwind", bson.D{{"path", "$lookup"}}}},
		{{"$match", bson.D{{"name", moveName}}}},
	}

	cursor, aggErr := moveCollection.Aggregate(context.TODO(), pipeline)
	if aggErr != nil {
		log.Println(aggErr.Error())
		return bson.M{}
	}

	var results []bson.M
	// cursor.
	if cursorErr := cursor.All(context.TODO(), &results); cursorErr != nil {
		log.Println(cursorErr.Error())
		return bson.M{}
	}

	if len(results) == 0 {
		return bson.M{}
	}

	return results[0]

}
