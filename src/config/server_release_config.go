//go:build release

package config

var LOG_REQUEST_AND_RESPONSE = false

func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}
