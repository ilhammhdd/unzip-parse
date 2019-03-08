package api

import (
	"encoding/json"
	"net/http"

	apiAdapter "bitbucket.com/bbd/unzip-parse/adapters/api"
)

func UnzipParseHandler(w http.ResponseWriter, r *http.Request) {

	apiAdapter.UnzipParse()

	jsonResponse, _ := json.Marshal(apiAdapter.Response{
		Message: "processed",
		Success: true,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
