package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	database "example.com/m/v2/database"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type ErrorMessage struct {
	Error string `json:"error"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user database.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := database.Db.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	response := map[string]interface{}{
		"message": "User created successfully",
		"user":    user,
	}

	json.NewEncoder(w).Encode(response)
}

func AddBalanceToUser(w http.ResponseWriter, r *http.Request) {
	var user database.User

	var requestPayload struct {
		Id      int     `json:"id"`
		Balance float64 `json:"balance"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if requestPayload.Id == 0 || requestPayload.Balance == 0 {
		http.Error(w, "Id or Balance is Required", http.StatusBadRequest)

		return
	}

	if err := database.Db.First(&user, requestPayload.Id).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "User doesn't exists"})

		return
	}

	newBalance := requestPayload.Balance + user.Balance
	if err := database.Db.Model(&user).Update("balance", newBalance).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	fmt.Fprintf(w, "User %d balance updated successfully to %f", user.Id, newBalance)
}

func GetBalanceOfUser(w http.ResponseWriter, r *http.Request) {
	var user database.User
	email := r.FormValue("email")

	if err := database.Db.Where("email = ?", email).First(&user).Error; err != nil {
		http.Error(w, "User doesn't exists", http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)

		return
	}

	if err := database.Db.First(&user, id).Error; err != nil {
		http.Error(w, "User doesn't exists", http.StatusBadRequest)

		return
	}

	fmt.Fprintf(w, "User %d has %f balance", user.Id, user.Balance)
}

func GetAllUsersBalance(w http.ResponseWriter, r *http.Request) {
	var users []database.User

	if err := database.Db.Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	type UserBalance struct {
		Id      int     `json:"id"`
		Balance float64 `json:"balance"`
	}

	var userBalances []UserBalance

	for _, user := range users {
		userBalances = append(userBalances, UserBalance{
			Id:      user.Id,
			Balance: user.Balance,
		})
	}

	response, err := json.Marshal(userBalances)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
