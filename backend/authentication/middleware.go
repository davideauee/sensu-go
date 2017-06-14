package authentication

import (
	"net/http"

	"github.com/sensu/sensu-go/backend/authentication/jwt"
)

// Middleware is a HTTP middleware that enforces authentication
func Middleware(provider Provider, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if authentication is enabled
		if !provider.AuthEnabled() {
			next.ServeHTTP(w, r)
			return
		}

		// TODO (Simon): We should probably avoid applying this middleware to these
		// routes instead...
		if r.URL.Path == "/auth" || r.URL.Path == "/auth/token" {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := jwt.ExtractBearerToken(r)
		if tokenString != "" {
			token, err := jwt.ValidateToken(tokenString)
			if err == nil {
				// Set the claims into the request context
				jwt.SetClaimsIntoContext(r, token)

				next.ServeHTTP(w, r)
				return
			}
		}

		// The user is not authenticated
		http.Error(w, "Request unauthorized", http.StatusUnauthorized)
		return
	})
}