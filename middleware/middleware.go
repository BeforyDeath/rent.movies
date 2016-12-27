package middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/BeforyDeath/rent.movies/handler"
)

func JSONContentType(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(fn)
}

func FetchJSONRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		// todo нуно приводить ключи к нижнему регистру!
		var req map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			if err != io.EOF {
				res := handler.Result{Error: "Invalid JSON parsing"}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(res)
				return
			}
		}

		ctx := context.WithValue(r.Context(), "request", req)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(fn)
}
