package db

import (
	"context"
	"go-maps/src"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Connect(collectionName string) *mongo.Collection {
	clientOptions := options.Client().ApplyURI(src.MongoURI)
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(context.TODO(), err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return client.Database(src.MongoDatabaseName).Collection(collectionName)
}

func Disconnect() {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}

func CreateIndexes(collection *mongo.Collection, indexNameParam string) {
	indexModel := []mongo.IndexModel{
		{
			Keys:    bson.M{indexNameParam: 1},
			Options: options.Index().SetUnique(true),
		},
	}

	indexView := collection.Indexes()
	cursor, err := indexView.List(context.TODO())
	if err != nil {
		log.Fatalf("Failed to list indexes: %v", err)
	}
	defer cursor.Close(context.TODO())

	existingIndexes := make(map[string]bool)
	for cursor.Next(context.TODO()) {
		var index bson.M
		if err := cursor.Decode(&index); err != nil {
			log.Fatalf("Failed to decode index: %v", err)
		}
		if name, ok := index["name"].(string); ok {
			existingIndexes[name] = true
		}
	}

	for _, model := range indexModel {
		indexName := model.Options.Name
		if indexName == nil {
			keys := model.Keys.(bson.M)
			for key := range keys {
				indexName = &key
				break
			}
		}
		if indexName != nil && !existingIndexes[*indexName] {
			_, err := indexView.CreateOne(context.TODO(), model)
			if err != nil {
				log.Fatalf("Failed to create index: %v", err)
			}
		}
	}
}
