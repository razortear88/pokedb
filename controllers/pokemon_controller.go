package controllers

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"time"

	"log"
	"strconv"
	"strings"

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

		file, formFileErr := c.FormFile("image")
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
			c.JSON(http.StatusBadRequest, responses.GameResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Only accept .png, .jpg, .jpeg extension  for image"}})
			return
		}

		f, openErr := file.Open()

		if openErr != nil {
			log.Fatal(openErr)
			return
		}

		cfg, configErr := config.LoadDefaultConfig(ctx)
		if configErr != nil {
			log.Printf("error: %v", configErr)
			return
		}

		// Create an Amazon S3 service client
		client := s3.NewFromConfig(cfg)
		uploader := manager.NewUploader(client)

		if len(c.Request.PostForm["name"]) != 0 {
			pokemon.Name = c.Request.PostForm["name"][0]
		}

		if len(c.Request.PostForm["nationalno"]) != 0 {
			nationalNo, err := strconv.Atoi(c.Request.PostForm["nationalno"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.NationalNo = nationalNo
		}

		if len(c.Request.PostForm["type[]"]) != 0 {
			if len(c.Request.PostForm["type[]"]) > 2 {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Pokemon Can Only Have 2 Types"}})
				return
			}
			pokemon.Type = c.Request.PostForm["type[]"]
		}

		if len(c.Request.PostForm["species"]) != 0 {
			pokemon.Species = c.Request.PostForm["species"][0]
		}

		if len(c.Request.PostForm["height"]) != 0 {
			height, heightErr := strconv.ParseFloat(c.Request.PostForm["height"][0], 32)
			if heightErr != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": heightErr.Error()}})
				return
			}
			pokemon.Height = float32(height)
		}

		if len(c.Request.PostForm["weight"]) != 0 {
			weight, weightErr := strconv.ParseFloat(c.Request.PostForm["weight"][0], 32)
			if weightErr != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": weightErr.Error()}})
				return
			}
			pokemon.Weight = float32(weight)
		}

		if len(c.Request.PostForm["basehp"]) != 0 {
			baseHp, err := strconv.Atoi(c.Request.PostForm["basehp"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseHP = baseHp
		}

		if len(c.Request.PostForm["minhp"]) != 0 {
			minHp, err := strconv.Atoi(c.Request.PostForm["minhp"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinHP = minHp
		}

		if len(c.Request.PostForm["maxhp"]) != 0 {
			maxHp, err := strconv.Atoi(c.Request.PostForm["maxhp"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxHP = maxHp
		}

		if len(c.Request.PostForm["baseattack"]) != 0 {
			baseAttack, err := strconv.Atoi(c.Request.PostForm["baseattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseAttack = baseAttack
		}

		if len(c.Request.PostForm["minattack"]) != 0 {
			minAttack, err := strconv.Atoi(c.Request.PostForm["minattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinAttack = minAttack
		}

		if len(c.Request.PostForm["maxattack"]) != 0 {
			maxAttack, err := strconv.Atoi(c.Request.PostForm["maxattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxAttack = maxAttack
		}

		if len(c.Request.PostForm["basedefense"]) != 0 {
			baseDefense, err := strconv.Atoi(c.Request.PostForm["basedefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseDefense = baseDefense
		}

		if len(c.Request.PostForm["mindefense"]) != 0 {
			minDefense, err := strconv.Atoi(c.Request.PostForm["mindefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinDefense = minDefense
		}

		if len(c.Request.PostForm["maxdefense"]) != 0 {
			maxDefense, err := strconv.Atoi(c.Request.PostForm["maxdefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxDefense = maxDefense
		}

		if len(c.Request.PostForm["basespattack"]) != 0 {
			baseSpAttack, err := strconv.Atoi(c.Request.PostForm["basespattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseSPAttack = baseSpAttack
		}

		if len(c.Request.PostForm["minspattack"]) != 0 {
			minSpAttack, err := strconv.Atoi(c.Request.PostForm["minspattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinSPAttack = minSpAttack
		}

		if len(c.Request.PostForm["maxspattack"]) != 0 {
			maxSpAttack, err := strconv.Atoi(c.Request.PostForm["maxspattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxSPAttack = maxSpAttack
		}

		if len(c.Request.PostForm["basespdefense"]) != 0 {
			baseSpDefense, err := strconv.Atoi(c.Request.PostForm["basespdefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseSPDefense = baseSpDefense
		}

		if len(c.Request.PostForm["minspdefense"]) != 0 {
			minSpDefense, err := strconv.Atoi(c.Request.PostForm["minspdefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinSPDefense = minSpDefense
		}

		if len(c.Request.PostForm["maxspdefense"]) != 0 {
			maxSpDefense, err := strconv.Atoi(c.Request.PostForm["maxspdefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxSPDefense = maxSpDefense
		}

		if len(c.Request.PostForm["basespeed"]) != 0 {
			baseSpeed, err := strconv.Atoi(c.Request.PostForm["basespeed"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseSpeed = baseSpeed
		}

		if len(c.Request.PostForm["minspeed"]) != 0 {
			minSpeed, err := strconv.Atoi(c.Request.PostForm["minspeed"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinSpeed = minSpeed
		}

		if len(c.Request.PostForm["maxspeed"]) != 0 {
			maxSpeed, err := strconv.Atoi(c.Request.PostForm["maxspeed"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxSpeed = maxSpeed
		}

		if len(c.Request.PostForm["total"]) != 0 {
			total, err := strconv.Atoi(c.Request.PostForm["total"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.Total = total
		}

		if validationErr := validate.Struct(&pokemon); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		var pokemonExist models.Pokemon

		err := pokemonCollection.FindOne(ctx, bson.M{"name": pokemon.Name}).Decode(&pokemonExist)
		if err == nil {
			c.JSON(http.StatusInternalServerError, responses.PokemonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Pokemon Already Exist"}})
			return
		}

		newPokemon := models.Pokemon{
			Name:          pokemon.Name,
			NationalNo:    pokemon.NationalNo,
			Type:          pokemon.Type,
			Species:       pokemon.Species,
			Height:        pokemon.Height,
			Weight:        pokemon.Weight,
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
		}

		pokemonCollection.InsertOne(ctx, newPokemon)

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

		update := bson.M{"image": result.Location}
		_, updateErr := pokemonCollection.UpdateOne(ctx, bson.M{"name": pokemon.Name}, bson.M{"$set": update})

		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, responses.PokemonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.Redirect(http.StatusFound, "/pokemon")

	}
}

func GetAllPokemons() []bson.M {

	pipeline := mongo.Pipeline{
		{{"$lookup", bson.D{{"from", "types"}, {"localField", "type"}, {"foreignField", "name"}, {"as", "lookup"}}}},
		{{"$unwind", bson.D{{"path", "$lookup"}}}},
		{{"$group", bson.D{{"_id", "$_id"}, {"name", bson.D{{"$first", "$name"}}}, {"lookup", bson.D{{"$push", "$lookup"}}}, {"nationalno", bson.D{{"$first", "$nationalno"}}}}}},
		{{"$sort", bson.D{{"nationalno", 1}}}},
	}

	cursor, aggErr := pokemonCollection.Aggregate(context.TODO(), pipeline)
	if aggErr != nil {
		log.Println(aggErr.Error())
		return []bson.M{}
	}

	var results []bson.M
	// cursor.
	if cursorErr := cursor.All(context.TODO(), &results); cursorErr != nil {
		log.Println(cursorErr.Error())
		return []bson.M{}
	}

	if len(results) == 0 {
		return []bson.M{}
	}

	return results

}

func GetPokemon(ctx *gin.Context) models.Pokemon {
	pokemonName := ctx.Param("pokemonName")
	var pokemon models.Pokemon

	pokemonCollection.FindOne(ctx, bson.M{"name": pokemonName}).Decode(&pokemon)

	return pokemon

}

func DeletePokemon() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		pokemonName := c.Param("pokemonName")
		defer cancel()

		var pokemon models.Pokemon

		filter := bson.D{{"name", pokemonName}}
		project := bson.D{{"image", 1}}
		opts := options.FindOne().SetProjection(project)
		dbErr := pokemonCollection.FindOne(ctx, filter, opts).Decode(&pokemon)
		if dbErr != nil {
			c.JSON(http.StatusInternalServerError, responses.PokemonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": dbErr.Error()}})
			return
		}

		result, err := pokemonCollection.DeleteOne(ctx, bson.M{"name": pokemonName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.PokemonResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.PokemonResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Pokemon with specified Name not found!"}},
			)
			return
		}

		imageFileName := strings.Split(pokemon.Image, "/")

		parsedUrl, parseErr := url.PathUnescape(imageFileName[len(imageFileName)-1])

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

		c.Redirect(http.StatusFound, "/pokemon")
	}
}

func GetDetailedPokemon(ctx *gin.Context) bson.M {
	pokemonName := ctx.Param("pokemonName")

	pipeline := mongo.Pipeline{
		{{"$lookup", bson.D{{"from", "types"}, {"localField", "type"}, {"foreignField", "name"}, {"as", "lookup"}}}},
		{{"$unwind", bson.D{{"path", "$lookup"}}}},
		{{"$match", bson.D{{"name", pokemonName}}}},
		{{"$group", bson.D{{"_id", "$_id"}, {"lookup", bson.D{{"$push", "$lookup"}}},
			{"name", bson.D{{"$first", "$name"}}}, {"nationalno", bson.D{{"$first", "$nationalno"}}}, {"species", bson.D{{"$first", "$species"}}},
			{"height", bson.D{{"$first", "$height"}}}, {"weight", bson.D{{"$first", "$weight"}}}, {"image", bson.D{{"$first", "$image"}}},
			{"basehp", bson.D{{"$first", "$basehp"}}}, {"minhp", bson.D{{"$first", "$minhp"}}}, {"maxhp", bson.D{{"$first", "$maxhp"}}},
			{"baseattack", bson.D{{"$first", "$baseattack"}}}, {"minattack", bson.D{{"$first", "$minattack"}}}, {"maxattack", bson.D{{"$first", "$maxattack"}}},
			{"basedefense", bson.D{{"$first", "$basedefense"}}}, {"mindefense", bson.D{{"$first", "$mindefense"}}}, {"maxdefense", bson.D{{"$first", "$maxdefense"}}},
			{"basespattack", bson.D{{"$first", "$basespattack"}}}, {"minspattack", bson.D{{"$first", "$minspattack"}}}, {"maxspattack", bson.D{{"$first", "$maxspattack"}}},
			{"basespdefense", bson.D{{"$first", "$basespdefense"}}}, {"minspdefense", bson.D{{"$first", "$minspdefense"}}}, {"maxspdefense", bson.D{{"$first", "$maxspdefense"}}},
			{"basespeed", bson.D{{"$first", "$basespeed"}}}, {"minspeed", bson.D{{"$first", "$minspeed"}}}, {"maxspeed", bson.D{{"$first", "$maxspeed"}}},
			{"total", bson.D{{"$first", "$total"}}}, {"prevevo", bson.D{{"$first", "$prevevo"}}}, {"nextevo", bson.D{{"$first", "$nextevo"}}},
		}}}}

	cursor, aggErr := pokemonCollection.Aggregate(context.TODO(), pipeline)
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

func EditPokemon() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		pokemonName := c.Param("pokemonName")
		var pokemon models.Pokemon

		defer cancel()

		file, _ := c.FormFile("image")

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
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Only accept .png, .jpg, .jpeg extension  for image"}})
				return
			}

			_, openErr := file.Open()
			if openErr != nil {
				log.Fatal(openErr)
				return
			}

		}

		if len(c.Request.PostForm["name"]) != 0 {
			pokemon.Name = c.Request.PostForm["name"][0]
		}

		if len(c.Request.PostForm["nationalno"]) != 0 {
			nationalNo, err := strconv.Atoi(c.Request.PostForm["nationalno"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.NationalNo = nationalNo
		}

		if len(c.Request.PostForm["type[]"]) != 0 {
			if len(c.Request.PostForm["type[]"]) > 2 {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Pokemon Can Only Have 2 Types"}})
				return
			}
			pokemon.Type = c.Request.PostForm["type[]"]
		}

		if len(c.Request.PostForm["species"]) != 0 {
			pokemon.Species = c.Request.PostForm["species"][0]
		}

		if len(c.Request.PostForm["height"]) != 0 {
			height, heightErr := strconv.ParseFloat(c.Request.PostForm["height"][0], 32)
			if heightErr != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": heightErr.Error()}})
				return
			}
			pokemon.Height = float32(height)
		}

		if len(c.Request.PostForm["weight"]) != 0 {
			weight, weightErr := strconv.ParseFloat(c.Request.PostForm["weight"][0], 32)
			if weightErr != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": weightErr.Error()}})
				return
			}
			pokemon.Weight = float32(weight)
		}

		if len(c.Request.PostForm["basehp"]) != 0 {
			baseHp, err := strconv.Atoi(c.Request.PostForm["basehp"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseHP = baseHp
		}

		if len(c.Request.PostForm["minhp"]) != 0 {
			minHp, err := strconv.Atoi(c.Request.PostForm["minhp"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinHP = minHp
		}

		if len(c.Request.PostForm["maxhp"]) != 0 {
			maxHp, err := strconv.Atoi(c.Request.PostForm["maxhp"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxHP = maxHp
		}

		if len(c.Request.PostForm["baseattack"]) != 0 {
			baseAttack, err := strconv.Atoi(c.Request.PostForm["baseattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseAttack = baseAttack
		}

		if len(c.Request.PostForm["minattack"]) != 0 {
			minAttack, err := strconv.Atoi(c.Request.PostForm["minattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinAttack = minAttack
		}

		if len(c.Request.PostForm["maxattack"]) != 0 {
			maxAttack, err := strconv.Atoi(c.Request.PostForm["maxattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxAttack = maxAttack
		}

		if len(c.Request.PostForm["basedefense"]) != 0 {
			baseDefense, err := strconv.Atoi(c.Request.PostForm["basedefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseDefense = baseDefense
		}

		if len(c.Request.PostForm["mindefense"]) != 0 {
			minDefense, err := strconv.Atoi(c.Request.PostForm["mindefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinDefense = minDefense
		}

		if len(c.Request.PostForm["maxdefense"]) != 0 {
			maxDefense, err := strconv.Atoi(c.Request.PostForm["maxdefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxDefense = maxDefense
		}

		if len(c.Request.PostForm["basespattack"]) != 0 {
			baseSpAttack, err := strconv.Atoi(c.Request.PostForm["basespattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseSPAttack = baseSpAttack
		}

		if len(c.Request.PostForm["minspattack"]) != 0 {
			minSpAttack, err := strconv.Atoi(c.Request.PostForm["minspattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinSPAttack = minSpAttack
		}

		if len(c.Request.PostForm["maxspattack"]) != 0 {
			maxSpAttack, err := strconv.Atoi(c.Request.PostForm["maxspattack"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxSPAttack = maxSpAttack
		}

		if len(c.Request.PostForm["basespdefense"]) != 0 {
			baseSpDefense, err := strconv.Atoi(c.Request.PostForm["basespdefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseSPDefense = baseSpDefense
		}

		if len(c.Request.PostForm["minspdefense"]) != 0 {
			minSpDefense, err := strconv.Atoi(c.Request.PostForm["minspdefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinSPDefense = minSpDefense
		}

		if len(c.Request.PostForm["maxspdefense"]) != 0 {
			maxSpDefense, err := strconv.Atoi(c.Request.PostForm["maxspdefense"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxSPDefense = maxSpDefense
		}

		if len(c.Request.PostForm["basespeed"]) != 0 {
			baseSpeed, err := strconv.Atoi(c.Request.PostForm["basespeed"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.BaseSpeed = baseSpeed
		}

		if len(c.Request.PostForm["minspeed"]) != 0 {
			minSpeed, err := strconv.Atoi(c.Request.PostForm["minspeed"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MinSpeed = minSpeed
		}

		if len(c.Request.PostForm["maxspeed"]) != 0 {
			maxSpeed, err := strconv.Atoi(c.Request.PostForm["maxspeed"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.MaxSpeed = maxSpeed
		}

		if len(c.Request.PostForm["total"]) != 0 {
			total, err := strconv.Atoi(c.Request.PostForm["total"][0])
			if err != nil {
				c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			pokemon.Total = total
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&pokemon); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.PokemonResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"name": pokemon.Name,
			"nationalno":    pokemon.NationalNo,
			"type":          pokemon.Type,
			"species":       pokemon.Species,
			"height":        pokemon.Height,
			"weight":        pokemon.Weight,
			"basehp":        pokemon.BaseHP,
			"minhp":         pokemon.MinHP,
			"maxhp":         pokemon.MaxHP,
			"baseattack":    pokemon.BaseAttack,
			"minattack":     pokemon.MinAttack,
			"maxattack":     pokemon.MaxAttack,
			"basedefense":   pokemon.BaseDefense,
			"mindefense":    pokemon.MinDefense,
			"maxdefense":    pokemon.MaxDefense,
			"basespattack":  pokemon.BaseSPAttack,
			"minspattack":   pokemon.MinSPAttack,
			"maxspattack":   pokemon.MaxSPAttack,
			"basespdefense": pokemon.BaseSPDefense,
			"minspdefense":  pokemon.MinSPDefense,
			"maxspdefense":  pokemon.MaxSPDefense,
			"basespeed":     pokemon.BaseSpeed,
			"minspeed":      pokemon.MinSpeed,
			"maxspeed":      pokemon.MaxSpeed,
			"total":         pokemon.Total,
		}
		_, err := pokemonCollection.UpdateOne(ctx, bson.M{"name": pokemonName}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.GameResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if file != nil {
			var oldImage models.Pokemon
			filter := bson.D{{"name", pokemon.Name}}
			project := bson.D{{"image", 1}}
			opts := options.FindOne().SetProjection(project)
			imageErr := pokemonCollection.FindOne(ctx, filter, opts).Decode(&oldImage)
			if imageErr != nil {
				c.JSON(http.StatusInternalServerError, responses.GameResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": imageErr.Error()}})
				return
			}
			// delete old image
			f, _ := file.Open()
			imageFileName := strings.Split(oldImage.Image, "/")

			parsedUrl, parseErr := url.PathUnescape((imageFileName[len(imageFileName)-1]))
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
			// upload new image
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

			update := bson.M{"image": result.Location}
			_, err := pokemonCollection.UpdateOne(ctx, bson.M{"name": pokemon.Name}, bson.M{"$set": update})
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.GameResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.Redirect(http.StatusFound, "/pokemon")

	}
}

func GetApiPokemons(ctx *gin.Context) {
	pipeline := mongo.Pipeline{
		{{"$project", bson.D{{"name", 1}, {"image", 1}, {"_id", 0}}}},
	}

	cursor, aggErr := pokemonCollection.Aggregate(context.TODO(), pipeline)
	if aggErr != nil {
		ctx.JSON(500, gin.H{"message": "Wrong Aggregate Pipeline"})
		return
	}

	var results []bson.M
	// cursor.
	if cursorErr := cursor.All(context.TODO(), &results); cursorErr != nil {
		ctx.JSON(500, gin.H{"message": "Failed to Get Data"})
		return
	}

	if len(results) == 0 {
		ctx.JSON(400, gin.H{"message": "No Pokemon Found"})
	}

	ctx.JSON(200, gin.H{
		"pokemons": results,
	})

}
