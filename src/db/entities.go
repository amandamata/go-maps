package db

type Address struct {
	Zipcode      string `json:"zipcode" bson:"zipcode"`
	Country      string `json:"country" bson:"country"`
	State        string `json:"state" bson:"state"`
	City         string `json:"city" bson:"city"`
	Neighborhood string `json:"neighborhood" bson:"neighborhood"`
	Street       string `json:"street" bson:"street"`
}
