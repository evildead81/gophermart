package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/evildead81/gophermart/internal/consts"
	"github.com/evildead81/gophermart/internal/helpers"
	"github.com/evildead81/gophermart/internal/storages"
)

func CreateOrder(storage storages.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(consts.UserIDKey).(int)

		var orderNumber string
		if err := json.NewDecoder(r.Body).Decode(&orderNumber); err != nil {
			http.Error(rw, "Invalid order number format", http.StatusBadRequest)
			return
		}

		if helpers.IsValidLuhn(orderNumber) {
			http.Error(rw, "Invalid order number format", http.StatusUnprocessableEntity)
			return
		}

		existUserID, err := storage.GetUserIdByOrderNumber(orderNumber)
		if err != nil && err != sql.ErrNoRows {
			http.Error(rw, "Server error", http.StatusInternalServerError)
			return
		}

		if err == nil && existUserID != int64(userID) {
			http.Error(rw, "Order number already used by another user", http.StatusConflict)
			return
		}

		if existUserID == int64(userID) {
			rw.WriteHeader(http.StatusAccepted)
			return
		}

		err = storage.CreateOrder(userID, orderNumber)
		if err != nil {
			http.Error(rw, "Server error", http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)

	}
}
