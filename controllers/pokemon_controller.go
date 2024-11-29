package controllers

import (
	"context"
	"net/http"
	"time"

	"log"
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

var pokemonCollection *mongo.Collection = configs.GetCollection(configs.DB, "pokemons")

func CreatePokemon() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pokemon models.Pokemon
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		//validate the request body
		if err := c.ShouldBind(&pokemon); err != nil {
			c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if len(c.Request.PostForm["name"]) != 0 {
			pokemon.Name = c.Request.PostForm["name"][0]
		}

		if len(c.Request.PostForm["nationalno"]) != 0 {
			pokemon.NationalNo = c.Request.PostForm["nationalno"][0]
		}

		log.Println(c.Request.PostForm["type[]"])
		if len(c.Request.PostForm["type[]"]) != 0 {
			pokemon.Type = c.Request.PostForm["type[]"]
		}

		if len(c.Request.PostForm["species"]) != 0 {
			pokemon.Species = c.Request.PostForm["species"][0]
		}

		if len(c.Request.PostForm["height"]) != 0 {
			height, heightErr := strconv.ParseFloat(c.Request.PostForm["height"][0], 32)
			if heightErr != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": heightErr.Error()}})
				return
			}
			pokemon.Height = float32(height)
		}

		if len(c.Request.PostForm["weight"]) != 0 {
			weight, weightErr := strconv.ParseFloat(c.Request.PostForm["weight"][0], 32)
			if weightErr != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": weightErr.Error()}})
				return
			}
			pokemon.Weight = float32(weight)
		}

		cfg, configErr := config.LoadDefaultConfig(ctx)
		if configErr != nil {
			log.Printf("error: %v", configErr)
			return
		}

		file, formFileErr := c.FormFile("image")
		if formFileErr != nil {
			log.Printf("error upload file")
			return
		}
		f, openErr := file.Open()

		if openErr != nil {
			log.Fatal(openErr)
			return
		}

		// Create an Amazon S3 service client
		client := s3.NewFromConfig(cfg)
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

		pokemon.Image = result.Location

		if len(c.Request.PostForm["basehp"]) != 0 {
			baseHp, err := strconv.Atoi(c.Request.PostForm["basehp"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseHP = baseHp
		}

		if len(c.Request.PostForm["minhp"]) != 0 {
			minHp, err := strconv.Atoi(c.Request.PostForm["minhp"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinHP = minHp
		}

		if len(c.Request.PostForm["maxhp"]) != 0 {
			maxHp, err := strconv.Atoi(c.Request.PostForm["maxhp"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxHP = maxHp
		}

		if len(c.Request.PostForm["baseattack"]) != 0 {
			baseAttack, err := strconv.Atoi(c.Request.PostForm["baseattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseAttack = baseAttack
		}

		if len(c.Request.PostForm["minattack"]) != 0 {
			minAttack, err := strconv.Atoi(c.Request.PostForm["minattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinAttack = minAttack
		}

		if len(c.Request.PostForm["maxattack"]) != 0 {
			maxAttack, err := strconv.Atoi(c.Request.PostForm["maxattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxAttack = maxAttack
		}

		if len(c.Request.PostForm["basedefense"]) != 0 {
			baseDefense, err := strconv.Atoi(c.Request.PostForm["basedefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseDefense = baseDefense
		}

		if len(c.Request.PostForm["mindefense"]) != 0 {
			minDefense, err := strconv.Atoi(c.Request.PostForm["mindefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinDefense = minDefense
		}

		if len(c.Request.PostForm["maxdefense"]) != 0 {
			maxDefense, err := strconv.Atoi(c.Request.PostForm["maxdefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxDefense = maxDefense
		}

		if len(c.Request.PostForm["basespattack"]) != 0 {
			baseSpAttack, err := strconv.Atoi(c.Request.PostForm["basespattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseSPAttack = baseSpAttack
		}

		if len(c.Request.PostForm["minspattack"]) != 0 {
			minSpAttack, err := strconv.Atoi(c.Request.PostForm["minspattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinSPAttack = minSpAttack
		}

		if len(c.Request.PostForm["maxspattack"]) != 0 {
			maxSpAttack, err := strconv.Atoi(c.Request.PostForm["maxspattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxSPAttack = maxSpAttack
		}

		if len(c.Request.PostForm["basespdefense"]) != 0 {
			baseSpDefense, err := strconv.Atoi(c.Request.PostForm["basespdefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseSPDefense = baseSpDefense
		}

		if len(c.Request.PostForm["minspdefense"]) != 0 {
			minSpDefense, err := strconv.Atoi(c.Request.PostForm["minspdefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinSPDefense = minSpDefense
		}

		if len(c.Request.PostForm["maxspdefense"]) != 0 {
			maxSpDefense, err := strconv.Atoi(c.Request.PostForm["maxspdefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxSPDefense = maxSpDefense
		}

		if len(c.Request.PostForm["basespeed"]) != 0 {
			baseSpeed, err := strconv.Atoi(c.Request.PostForm["basespeed"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseSpeed = baseSpeed
		}

		if len(c.Request.PostForm["minspeed"]) != 0 {
			minSpeed, err := strconv.Atoi(c.Request.PostForm["minspeed"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinSpeed = minSpeed
		}

		if len(c.Request.PostForm["maxspeed"]) != 0 {
			maxSpeed, err := strconv.Atoi(c.Request.PostForm["maxspeed"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxSpeed = maxSpeed
		}

		if len(c.Request.PostForm["total"]) != 0 {
			total, err := strconv.Atoi(c.Request.PostForm["total"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.Total = total
		}

		if len(c.Request.PostForm["prevevo"]) != 0 {
			pokemon.PrevEvo = c.Request.PostForm["prevevo"][0]
		}

		if len(c.Request.PostForm["nextevo"]) != 0 {
			pokemon.NextEvo = c.Request.PostForm["nextevo"][0]
		}

		if validationErr := validate.Struct(&pokemon); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		var pokemonExist models.Pokemon

		err := pokemonCollection.FindOne(ctx, bson.M{"name": pokemon.Name}).Decode(&pokemonExist)
		if err == nil {
			c.JSON(http.StatusInternalServerError, responses.GameResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Game Already Exist"}})
			return
		}

		newPokemon := models.Pokemon{
			Name:          pokemon.Name,
			NationalNo:    pokemon.NationalNo,
			Type:          pokemon.Type,
			Species:       pokemon.Species,
			Height:        pokemon.Height,
			Weight:        pokemon.Weight,
			Image:         pokemon.Image,
			BaseHP:        pokemon.BaseHP,
			MinHP:         pokemon.MinHP,
			MaxHP:         pokemon.MaxHP,
			BaseAttack:    pokemon.BaseAttack,
			MinAttack:     pokemon.MinAttack,
			MaxAttack:     pokemon.MaxAttack,
			BaseDefense:   pokemon.BaseDefense,
			MinDefense:    pokemon.MinDefense,
			MaxDefense:    pokemon.MaxDefense,
			BaseSPAttack:  pokemon.BaseSPAttack,
			MinSPAttack:   pokemon.MinSPAttack,
			MaxSPAttack:   pokemon.MaxSPAttack,
			BaseSPDefense: pokemon.BaseSPDefense,
			MinSPDefense:  pokemon.MinSPDefense,
			MaxSPDefense:  pokemon.MaxSPDefense,
			BaseSpeed:     pokemon.BaseSpeed,
			MinSpeed:      pokemon.MinSpeed,
			MaxSpeed:      pokemon.MaxSpeed,
			Total:         pokemon.Total,
			PrevEvo:       pokemon.PrevEvo,
			NextEvo:       pokemon.NextEvo,
		}

		pokemonCollection.InsertOne(ctx, newPokemon)

		c.Redirect(http.StatusFound, "/pokemon")

	}
}

func GetAllPokemons() []models.Pokemon {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var pokemons []models.Pokemon
	defer cancel()
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{"name", 1}})

	results, err := pokemonCollection.Find(ctx, filter, opts)

	if err != nil {
		return pokemons
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singlePokemon models.Pokemon
		if err = results.Decode(&singlePokemon); err != nil {
			return pokemons
		}
		pokemons = append(pokemons, singlePokemon)
	}

	return pokemons
}
