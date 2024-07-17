package main

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectMongo() (*mongo.Client, error) {
	mongoUri := os.Getenv("MONGO_URI")
	if mongoUri == "" {
		return nil, errors.New("MONGO_URI environment variable is not set")
	}

	username := os.Getenv("MONGO_USERNAME")
	if username == "" {
		return nil, errors.New("MONGO_USERNAME environment variable is not set")
	}

	password := os.Getenv("MONGO_PASSWORD")
	if password == "" {
		return nil, errors.New("MONGO_PASSWORD environment variable is not set")
	}

	clientOptions := options.Client().ApplyURI(mongoUri)
	clientOptions.SetAuth(options.Credential{
		Username: username,
		Password: password,
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = c.Ping(context.TODO(), nil)
	
	if err != nil {
		return nil, err
	}

	return c, nil
}
