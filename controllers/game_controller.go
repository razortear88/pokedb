package controllers

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/razortear88/pokedb/configs"
	"github.com/razortear88/pokedb/models"
	"github.com/razortear88/pokedb/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var gameCollection *mongo.Collection = configs.GetCollection(configs.DB, "games")

// var validate = validator.New()

func CreateGame() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var game models.Game
		defer cancel()

		file, formFileErr := c.FormFile("cover")

		if formFileErr != nil {
			log.Printf("error upload file")
			return
		}

		fileExtension := path.Ext(file.Filename) //obtain the extension of file
		rightFileExtension := false
		for _, extension := range [3]string{".png", ".jpeg", ".jpg"} {
			if fileExtension == extension {
				rightFileExtension = true
				break
			}
		}

		if !rightFileExtension {
			c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Only accept .png, .jpg, .jpeg extension  for cover"}})
			return
		}

		f, openErr := file.Open()

		if openErr != nil {
			log.Fatal(openErr)
			return
		}

		// setup s3 uploader
		cfg, configErr := config.LoadDefaultConfig(ctx)
		if configErr != nil {
			log.Printf("error: %v", configErr)
			return
		}

		// Create an Amazon S3 service client
		client := s3.NewFromConfig(cfg)
		uploader := manager.NewUploader(client)

		//validate the request body
		if err := c.ShouldBind(&game); err != nil {
			c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if len(c.Request.PostForm["name"]) != 0 {
			game.Name = c.Request.PostForm["name"][0]
		}
		if len(c.Request.PostForm["generation"]) != 0 {
			generation, err := strconv.Atoi(c.Request.PostForm["generation"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			game.Generation = generation
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&game); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		var gameExist models.Game

		err := gameCollection.FindOne(ctx, bson.M{"name": game.Name}).Decode(&gameExist)
		if err == nil {
			c.JSON(http.StatusInternalServerError, responses.GameResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Game Already Exist"}})
			return
		}

		newGame := models.Game{
			Name:       game.Name,
			Generation: game.Generation,
		}

		gameCollection.InsertOne(ctx, newGame)

		result, uploadErr := uploader.Upload(ctx, &s3.PutObjectInput{
			Bucket: aws.String("my-pokedb-project"),
			Key:    aws.String(file.Filename),
			Body:   f,
			ACL:    "public-read",
		})
		if uploadErr != nil {
			log.Fatal(uploadErr)
			return
		}

		update := bson.M{"cover": result.Location}
		_, updateErr := gameCollection.UpdateOne(ctx, bson.M{"name": game.Name}, bson.M{"$set": update})

		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, responses.GameResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.Redirect(http.StatusFound, "/game")
	}
}

func GetAllGames() []models.Game {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var games []models.Game
	defer cancel()
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{"name", 1}})

	results, err := gameCollection.Find(ctx, filter, opts)

	if err != nil {
		return games
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleGame models.Game
		if err = results.Decode(&singleGame); err != nil {
			return games
		}
		games = append(games, singleGame)
	}
	return games
}

func GetGame(ctx *gin.Context) models.Game {
	gameName := ctx.Param("gameName")
	var game models.Game

	gameCollection.FindOne(ctx, bson.M{"name": gameName}).Decode(&game)

	return game

}

func EditGame() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		gameName := c.Param("gameName")
		var game models.Game

		defer cancel()

		file, _ := c.FormFile("cover")

		if file != nil {
			fileExtension := path.Ext(file.Filename) //obtain the extension of file
			rightFileExtension := false
			for _, extension := range [3]string{".png", ".jpeg", ".jpg"} {
				if fileExtension == extension {
					rightFileExtension = true
					break
				}
			}

			if !rightFileExtension {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Only accept .png, .jpg, .jpeg extension  for cover"}})
				return
			}

			_, openErr := file.Open()
			if openErr != nil {
				log.Fatal(openErr)
				return
			}

		}

		if len(c.Request.PostForm["name"]) != 0 {
			game.Name = c.Request.PostForm["name"][0]
		}

		if len(c.Request.PostForm["generation"]) != 0 {
			generation, err := strconv.Atoi(c.Request.PostForm["generation"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			game.Generation = generation
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&game); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"name": game.Name, "generation": game.Generation}
		_, err := gameCollection.UpdateOne(ctx, bson.M{"name": gameName}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GameResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if file != nil {
			var oldCover models.Game
			filter := bson.D{{"name", game.Name}}
			project := bson.D{{"cover", 1}}
			opts := options.FindOne().SetProjection(project)
			coverErr := gameCollection.FindOne(ctx, filter, opts).Decode(&oldCover)
			if coverErr != nil {
				c.JSON(http.StatusInternalServerError, responses.GameResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": coverErr.Error()}})
				return
			}
			// delete old cover
			f, _ := file.Open()
			coverFileName := strings.Split(oldCover.Cover, "/")

			parsedUrl, parseErr := url.PathUnescape((coverFileName[len(coverFileName)-1]))
			if parseErr != nil {
				log.Printf("error: %v", parseErr)
				return
			}

			input := &s3.DeleteObjectInput{
				Bucket: aws.String("my-pokedb-project"),
				Key:    aws.String(parsedUrl),
			}
			cfg, configErr := config.LoadDefaultConfig(ctx)

			if configErr != nil {
				log.Printf("error: %v", configErr)
				return
			}

			// Create an Amazon S3 service client
			client := s3.NewFromConfig(cfg)
			_, DeleteErr := client.DeleteObject(c, input)

			if DeleteErr != nil {
				log.Printf("error: %v", DeleteErr)
				return
			}
			// upload new cover
			uploader := manager.NewUploader(client)

			result, uploadErr := uploader.Upload(ctx, &s3.PutObjectInput{
				Bucket: aws.String("my-pokedb-project"),
				Key:    aws.String(file.Filename),
				Body:   f,
				ACL:    "public-read",
			})

			if uploadErr != nil {
				log.Fatal(uploadErr)
				return
			}

			update := bson.M{"cover": result.Location}
			_, err := gameCollection.UpdateOne(ctx, bson.M{"name": game.Name}, bson.M{"$set": update})
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.GameResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.Redirect(http.StatusFound, "/game")
	}
}

func DeleteGame() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		gameName := c.Param("gameName")
		defer cancel()

		var game models.Game

		filter := bson.D{{"name", gameName}}
		project := bson.D{{"cover", 1}}
		opts := options.FindOne().SetProjection(project)
		dbErr := gameCollection.FindOne(ctx, filter, opts).Decode(&game)
		if dbErr != nil {
			c.JSON(http.StatusInternalServerError, responses.GameResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": dbErr.Error()}})
			return
		}

		result, err := gameCollection.DeleteOne(ctx, bson.M{"name": gameName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GameResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.GameResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Game with specified Name not found!"}},
			)
			return
		}

		coverFileName := strings.Split(game.Cover, "/")

		parsedUrl, parseErr := url.PathUnescape(coverFileName[len(coverFileName)-1])

		if parseErr != nil {
			log.Printf("error: %v", parseErr)
			return
		}

		input := &s3.DeleteObjectInput{
			Bucket: aws.String("my-pokedb-project"),
			Key:    aws.String(parsedUrl),
		}

		cfg, configErr := config.LoadDefaultConfig(ctx)

		if configErr != nil {
			log.Printf("error: %v", configErr)
			return
		}

		// Create an Amazon S3 service client
		client := s3.NewFromConfig(cfg)
		_, DeleteErr := client.DeleteObject(c, input)

		if DeleteErr != nil {
			log.Printf("error: %v", DeleteErr)
			return
		}

		c.Redirect(http.StatusFound, "/game")
	}
}
