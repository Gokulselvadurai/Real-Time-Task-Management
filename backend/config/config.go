package config

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	JWTSecret      []byte
	Client         *mongo.Client
	TaskCollection *mongo.Collection
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	JWTSecret = []byte(os.Getenv("JWT_SECRET"))

	mongoURI := os.Getenv("MONGODB_URI")
	Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	TaskCollection = Client.Database("taskdb").Collection("tasks")
}
