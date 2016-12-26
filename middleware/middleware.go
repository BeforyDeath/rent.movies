package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/BeforyDeath/rent.movies/handler"
)

func JsonContentType(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(fn)
}

func FetchJsonRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()
		var req map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			if err != io.EOF {
				res := handler.Result{Error: "invalid JSON Request"}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(res)
				// todo сделать логи
				fmt.Printf("Invalid JSON Request: %v\n", err)
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