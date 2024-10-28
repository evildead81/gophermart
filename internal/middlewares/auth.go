package middlewares

import (
	"context"
	"net/http"

	"github.com/evildead81/gophermart/internal/consts"
	"github.com/evildead81/gophermart/internal/session"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(consts.CookieName)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := session.ValidateAuthToken(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), consts.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
