package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/keslerliv/wallet/internal/entities"
)

func CreateWallet(w http.ResponseWriter, r *http.Request) {
	var wallet entities.Wallet
	var payload entities.CreateWalletPayload

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	walletData := wallet.NewWallet(wallet.UserID)

	response := map[string]any{"id": walletData.ID}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
