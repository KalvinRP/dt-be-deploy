package routes

import (
	"dewetour/2pkg/middleware"
	"dewetour/2pkg/mysql"
	repositories "dewetour/4repositories"
	handlers "dewetour/6handlers"

	"github.com/gorilla/mux"
)

func TransRoutes(r *mux.Router) {
	transRepository := repositories.RepositoryTrans(mysql.DB)
	h := handlers.HandlerTransaction(transRepository)

	r.HandleFunc("/transaction", h.FindTrans).Methods("GET")
	r.HandleFunc("/transaction/{id}", h.GetTrans).Methods("GET")
	r.HandleFunc("/account", middleware.Auth(h.UserHistory)).Methods("GET")
	r.HandleFunc("/transaction", middleware.Auth(h.MakeTrans)).Methods("POST")
	r.HandleFunc("/transaction/{id}", middleware.Auth(h.EditTrans)).Methods("PATCH")
	r.HandleFunc("/transaction/{id}", middleware.Auth(h.DeleteTrans)).Methods("DELETE")
	r.HandleFunc("/notification", h.Notification).Methods("POST")
}
