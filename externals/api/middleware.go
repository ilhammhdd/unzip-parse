package api

import (
	"encoding/json"
	"net/http"

	apiAdapter "bitbucket.com/bbd/unzip-parse/adapters/api"
)

func MustMethod(method string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			jsonResponse, _ := json.Marshal(apiAdapter.Response{
				Message: "method not allowed",
				Success: false,
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(jsonResponse)
			return
		}

		h.ServeHTTP(w, r)
	})
}
