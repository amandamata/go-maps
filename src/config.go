package src

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	MongoURI            string
	MongoDatabaseName   string
	MongoCollectionName string
	GoogleMapsApiKey    string
	GoogleMapsApiUrl    string
)

func LoadEnvs() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	GoogleMapsApiKey = os.Getenv("GOOGLE_MAPS_API_KEY")
	GoogleMapsApiUrl = os.Getenv("GOOGLE_MAPS_API_URL")

	if GoogleMapsApiKey == "" {
		log.Fatal("GOOGLE_MAPS_API_KEY is missing")
	}

	if GoogleMapsApiUrl == "" {
		log.Fatal("GOOGLE_MAPS_API_URL is missing")
	}
}

func LoadDBEnvs() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MongoURI = os.Getenv("MONGO_URI")
	MongoDatabaseName = os.Getenv("MONGO_DATABASE_NAME")
	MongoCollectionName = os.Getenv("MONGO_COLLECTION_NAME")

	if MongoURI == "" {
		log.Fatal("MONGO_URI is missing")
	}

	if MongoDatabaseName == "" {
		log.Fatal("MONGO_DATABASE_NAME is missing")
	}

	if MongoCollectionName == "" {
		log.Fatal("MONGO_COLLECTION_NAME is missing")
	}
}
