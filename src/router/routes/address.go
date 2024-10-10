package routes

import (
	"go-maps/src/controllers"
	"net/http"
)

var addressController = controllers.NewAddressController()

type Route struct {
	Uri                    string
	Method                 string
	Function               http.HandlerFunc
	AuthenticationRequired bool
}

var loginRoutes = []Route{
	{
		Uri:                    "/zipcode",
		Method:                 http.MethodGet,
		Function:               addressController.Zipcode,
		AuthenticationRequired: false,
	},
}
