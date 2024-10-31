package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/evildead81/gophermart/internal/consts"
	"github.com/evildead81/gophermart/internal/contracts"
	"github.com/evildead81/gophermart/internal/errors"
	"github.com/evildead81/gophermart/internal/storages"
)

func Withdraw(storage storages.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(consts.UserIDKey).(int)

		var request contracts.OrderRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(rw, "Invalid request format", http.StatusBadRequest)
			return
		}

		err := storage.Withdraw(userID, request.Order, request.Sum)
		if err == errors.ErrPaymentRequiredError {
			http.Error(rw, "Insufficient funds", http.StatusPaymentRequired)
			return
		} else if err != nil {
			http.Error(rw, "Server error", http.StatusInternalServerError)
		}

		rw.WriteHeader(http.StatusOK)
	}
}
