package router

import (
	"net/http"

	"github.com/gorilla/mux"

	handlers "example.com/m/v2/handlers"
)

func Start() {
	router := mux.NewRouter()

	router.HandleFunc("/create-user", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/add-balance-to-user", handlers.AddBalanceToUser).Methods("POST")
	router.HandleFunc("/get-balance-of-user/{id}", handlers.GetBalanceOfUser).Methods("GET")
	router.HandleFunc("/get-all-users-balance", handlers.GetAllUsersBalance).Methods("GET")
	router.HandleFunc("/transfer-balance", handlers.TransferBalance).Methods("POST")
	router.HandleFunc("/withdraw", handlers.Withdraw).Methods("POST")

	http.ListenAndServe(":8080", router)
}
