package routes

import (
	"dewetour/2pkg/middleware"
	"dewetour/2pkg/mysql"
	repositories "dewetour/4repositories"
	handlers "dewetour/6handlers"

	"github.com/gorilla/mux"
)

func TripsRoutes(r *mux.Router) {
	tripsRepository := repositories.RepositoryTrips(mysql.DB)
	h := handlers.HandlerTrips(tripsRepository)

	r.HandleFunc("/trips", h.FindTrips).Methods("GET")
	r.HandleFunc("/trips/{id}", h.GetTrips).Methods("GET")
	r.HandleFunc("/trips", middleware.AuthAdmin(middleware.UploadFile(h.MakeTrips))).Methods("POST")
	r.HandleFunc("/trips/{id}", middleware.AuthAdmin(middleware.UploadFile(h.EditTrips))).Methods("PATCH")
	r.HandleFunc("/trips/{id}", middleware.AuthAdmin(h.DeleteTrips)).Methods("DELETE")
}
