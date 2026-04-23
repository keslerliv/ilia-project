package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/keslerliv/user/internal/entities"
	"github.com/keslerliv/user/internal/models"
	"github.com/keslerliv/user/internal/services/kafka"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var payload entities.LoginPayload

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get user by email
	user, err := models.GetUserByEmail(payload.Email)
	if err != nil {
		http.Error(w, "invalid user or password", http.StatusInternalServerError)
		return
	}

	// validate user password
	if !user.ValidatePassword(payload.Password) {
		http.Error(w, "invalid user or password", http.StatusInternalServerError)
		return
	}

	// get user token
	token, err := user.CreateJWT()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{"token": token}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	var payload entities.CreateUserPayload

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userData, err := user.NewUser(payload.Email, payload.Password)

	user_id, err := models.CreateUser(userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// request wallet by kafka
	userMessage := map[string]interface{}{"action": "new_user", "user_id": user_id}
	data, err := json.Marshal(userMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = kafka.PushOrderToQueue([]byte(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{"id": user_id}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
