package routes

import (
	"github.com/gorilla/mux"
)

func RouteInit(r *mux.Router) {
	UserRoutes(r)
	TripsRoutes(r)
	CountryRoutes(r)
	AuthRoutes(r)
	TransRoutes(r)
}
