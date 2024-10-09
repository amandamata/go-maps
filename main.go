package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const googleMapsAPIURL = "https://maps.googleapis.com/maps/api/geocode/json"

type GeocodeResponse struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
	} `json:"results"`
	Status string `json:"status"`
}

type Address struct {
	CEP          string `json:"cep" bson:"cep"`
	Country      string `json:"country" bson:"country"`
	State        string `json:"state" bson:"state"`
	City         string `json:"city" bson:"city"`
	Neighborhood string `json:"neighborhood" bson:"neighborhood"`
	Street       string `json:"street" bson:"street"`
}

var bloomFilter *bloom.BloomFilter

func initBloomFilter() {
	bloomFilter = bloom.NewWithEstimates(10000, 0.01)
}

func getGeocodeData(cep, apiKey string) (*GeocodeResponse, error) {
	url := fmt.Sprintf("%s?address=%s&key=%s", googleMapsAPIURL, cep, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var geocodeResponse GeocodeResponse
	if err := json.Unmarshal(body, &geocodeResponse); err != nil {
		return nil, err
	}

	return &geocodeResponse, nil
}

func saveAddressToMongo(address Address) error {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return fmt.Errorf("MONGO_URI is missing")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("failed to create mongo client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to mongo: %v", err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("geocode").Collection("addresses")
	_, err = collection.InsertOne(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to insert address: %v", err)
	}

	bloomFilter.AddString(address.CEP)

	return nil
}

func getAddressFromMongo(cep string) (*Address, error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI is missing")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %v", err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("geocode").Collection("addresses")
	var address Address
	err = collection.FindOne(ctx, bson.M{"cep": cep}).Decode(&address)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find address: %v", err)
	}

	return &address, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")

	if cep == "" {
		http.Error(w, "CEP is required", http.StatusBadRequest)
		return
	}

	if apiKey == "" {
		http.Error(w, "API Key is missing", http.StatusInternalServerError)
		return
	}

	if bloomFilter.TestString(cep) {
		address, err := getAddressFromMongo(cep)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if address != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(address)
			return
		}
	}

	geocodeData, err := getGeocodeData(cep, apiKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if geocodeData.Status != "OK" {
		http.Error(w, "Failed to get geocode data", http.StatusInternalServerError)
		return
	}

	address := &Address{CEP: cep}
	for _, component := range geocodeData.Results[0].AddressComponents {
		for _, t := range component.Types {
			switch t {
			case "country":
				address.Country = component.LongName
			case "administrative_area_level_1":
				address.State = component.LongName
			case "administrative_area_level_2":
				address.City = component.LongName
			case "sublocality_level_1", "sublocality", "neighborhood":
				address.Neighborhood = component.LongName
			case "route":
				address.Street = component.LongName
			}
		}
	}

	err = saveAddressToMongo(*address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(address)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	initBloomFilter()

	http.HandleFunc("/geocode", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
