package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/keslerliv/user/internal/entities"
	"github.com/keslerliv/user/internal/models"
	"github.com/keslerliv/user/internal/services/kafka"
)

func PostUser(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	var payload entities.CreateUserPayload

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userData, err := user.NewUser(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userResponse, err := models.PostUser(userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// request wallet by kafka
	userMessage := map[string]interface{}{"action": "new_user", "user_id": userResponse.ID}
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResponse)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	user, err := models.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func PatchUser(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	var payload entities.CreateUserPayload

	idStr := chi.URLParam(r, "id")
	uid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Decode request body
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userData, err := user.NewPatchUser(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userResponse, err := models.PatchUser(uid, map[string]interface{}{
		"first_name": userData.FirstName,
		"last_name":  userData.LastName,
		"email":      userData.Email,
		"password":   userData.Password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResponse)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = models.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	var payload entities.LoginPayload

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get user by email
	user, err := models.GetUserByEmail(payload.User.Email)
	if err != nil {
		http.Error(w, "invalid user or password", http.StatusInternalServerError)
		return
	}

	// validate user password
	if !user.ValidatePassword(payload.User.Password) {
		http.Error(w, "invalid user or password", http.StatusInternalServerError)
		return
	}

	// get user token
	token, err := user.CreateJWT()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{"user": map[string]any{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
	}, "token": token}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
