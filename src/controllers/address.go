package controllers

import (
	"encoding/json"
	"fmt"
	"go-maps/src/db"
	"go-maps/src/models"
	"go-maps/src/repositories"
	"net/http"
	"os"

	"github.com/bits-and-blooms/bloom/v3"
)

type AddressController struct {
	repo repositories.AddressRepository
}

func NewAddressController() *AddressController {
	addressRepository, err := repositories.NewRepository()
	if err != nil {
		return nil
	}
	return &AddressController{repo: *addressRepository}
}

var bloomFilter *bloom.BloomFilter

func initBloomFilter() {
	bloomFilter = bloom.NewWithEstimates(10000, 0.01)
}

func (c *AddressController) Zipcode(w http.ResponseWriter, r *http.Request) {
	zipcode := r.URL.Query().Get("zipcode")

	if zipcode == "" {
		http.Error(w, "Zipcode is required", http.StatusBadRequest)
		return
	}

	if bloomFilter.TestString(zipcode) {
		address, err := c.repo.FindByZipcode(zipcode)
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

	geocodeData, err := getGeocodeData(zipcode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if geocodeData.Status != "OK" {
		http.Error(w, "Failed to get geocode data", http.StatusInternalServerError)
		return
	}

	address := &db.Address{Zipcode: zipcode}
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

	err = c.repo.Save(*address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bloomFilter.AddString(address.Zipcode)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(address)
}

func getGeocodeData(zipcode string) (*models.GeocodeResponse, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	apiUrl := os.Getenv("GOOGLE_MAPS_API_URL")

	url := fmt.Sprintf("%s?address=%s&key=%s", apiUrl, zipcode, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var geocodeResponse models.GeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&geocodeResponse); err != nil {
		return nil, err
	}

	return &geocodeResponse, nil
}
