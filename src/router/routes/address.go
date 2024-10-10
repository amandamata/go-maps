package routes

import (
	"go-maps/src/controllers"
	"net/http"
)

var addressController = controllers.NewAddressController()

var addressRoutes = []Route{
	{
		Uri:                    "/address",
		Method:                 http.MethodGet,
		Function:               addressController.Zipcode,
		AuthenticationRequired: false,
	},
}
