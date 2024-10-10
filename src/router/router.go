package router

import (
	"github.com/gorilla/mux"

	httpSwagger "github.com/swaggo/http-swagger"
)

func Generate() *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	return router
}
