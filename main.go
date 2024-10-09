package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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

func handler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if cep == "" {
		http.Error(w, "CEP is required", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		http.Error(w, "Google Maps API key is not set", http.StatusInternalServerError)
		return
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

	address := make(map[string]string)
	for _, component := range geocodeData.Results[0].AddressComponents {
		for _, t := range component.Types {
			switch t {
			case "country":
				address["country"] = component.LongName
			case "administrative_area_level_1":
				address["state"] = component.LongName
			case "administrative_area_level_2":
				address["city"] = component.LongName
			case "sublocality_level_1", "sublocality", "neighborhood":
				address["neighborhood"] = component.LongName
			case "route":
				address["street"] = component.LongName
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(address)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	http.HandleFunc("/geocode", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
