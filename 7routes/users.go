package routes

import (
	"dewetour/2pkg/middleware"
	"dewetour/2pkg/mysql"
	repositories "dewetour/4repositories"
	handlers "dewetour/6handlers"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	userRepository := repositories.RepositoryUser(mysql.DB)
	h := handlers.HandlerUser(userRepository)

	r.HandleFunc("/user", middleware.Auth(h.FindAcc)).Methods("GET")
	r.HandleFunc("/useracc", middleware.Auth(h.GetAcc)).Methods("GET")
	r.HandleFunc("/user", middleware.Auth(middleware.UploadFile(h.EditAcc))).Methods("PATCH")
	r.HandleFunc("/user/{id}", middleware.Auth(h.DeleteAcc)).Methods("DELETE")
}
