package supertokens

import (
	"net/http"
)

// Middleware for verifying and refreshing session
func Middleware(doAntiCsrfCheck ...bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO:
	})
}
