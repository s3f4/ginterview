package app

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoClient is used to connect mongodb, connection will be done one time.
var mongoClient *mongo.Client

// connectMongo connects mongodb with mongodb dsn string
func connectMongo(dsn string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(dsn)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	// Check the connection
	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	return client, nil
}

// initConnections starts all connections
func initConnections() {
	var err error
	mongoClient, err = connectMongo(os.Getenv("MONGO_DSN"))
	if err != nil {
		log.Panic(err)
	}
}
