package router

import (
	"github.com/gorilla/mux"

	"go-maps/src/router/routes"

	httpSwagger "github.com/swaggo/http-swagger"
)

func Generate() *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	return routes.Config(router)
}
