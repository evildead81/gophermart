package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/evildead81/gophermart/internal/contracts"
	"github.com/evildead81/gophermart/internal/errors"
	"github.com/evildead81/gophermart/internal/session"
	"github.com/evildead81/gophermart/internal/storages"
)

func Register(storage storages.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var request contracts.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(rw, "Invalid request format", http.StatusBadRequest)
			return
		}

		err := storage.RegisterUser(request.Login, request.Password)

		if err == errors.UserIsAlreadyExists {
			http.Error(rw, "Login is already taken", http.StatusConflict)
			return
		}

		if err != nil {
			http.Error(rw, "Internal server error", http.StatusInternalServerError)
			return
		}

		userID, err := storage.GetUserIdByLogin(request.Login)
		if err != nil {
			http.Error(rw, "Internal server error", http.StatusInternalServerError)
			return
		}

		token, err := session.GenerateAuthToken(int(userID))
		if err != nil {
			http.Error(rw, "Error generating token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(rw, &http.Cookie{
			Name:     "AuthToken",
			Value:    token,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})

		rw.WriteHeader(http.StatusOK)
	}
}
