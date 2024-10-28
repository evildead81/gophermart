package handlers

import (
	"net/http"

	"github.com/evildead81/gophermart/internal/storages"
)

func Login(storage storages.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

	}
}
