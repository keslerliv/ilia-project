package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/keslerliv/user/internal/entities"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	var payload entities.CreateUserPayload

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userData := user.NewUser(user.ID)

	response := map[string]any{"id": userData.ID}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
