package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(EnvMongoURI()).SetServerAPIOptions(serverAPI)
	client, dbErr := mongo.Connect(context.TODO(), opts)
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//ping the database
	dbErr = client.Ping(ctx, nil)
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

// Client instance
var DB *mongo.Client = ConnectDB()

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("golangAPI").Collection(collectionName)
	return collection
}
