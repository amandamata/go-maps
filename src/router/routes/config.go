package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Uri                    string
	Method                 string
	Function               func(http.ResponseWriter, *http.Request)
	AuthenticationRequired bool
}

func Config(router *mux.Router) *mux.Router {
	routes := addressRoutes
	for _, route := range routes {
		if route.AuthenticationRequired {
			router.HandleFunc(route.Uri, route.Function).Methods(route.Method, http.MethodOptions)
		} else {
			router.HandleFunc(route.Uri, route.Function).Methods(route.Method, http.MethodOptions)
		}
	}
	return router
}
