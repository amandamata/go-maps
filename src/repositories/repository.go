package repositories

import (
	"context"
	"fmt"
	"go-maps/src/db"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AddressRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewRepository() (*AddressRepository, error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI is missing")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %v", err)
	}

	collection := client.Database("geocode").Collection("addresses")

	return &AddressRepository{
		client:     client,
		collection: collection,
	}, nil
}

func (r *AddressRepository) Save(address db.Address) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to insert address: %v", err)
	}

	return nil
}

func (r *AddressRepository) FindByZipcode(zipcode string) (*db.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var address db.Address
	err := r.collection.FindOne(ctx, bson.M{"zipcode": zipcode}).Decode(&address)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find address: %v", err)
	}

	return &address, nil
}
