package api

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Data any `json:"data"`
}

// JSONResponse writes a JSON response with the given status code and data.
func JSONResponse(res http.ResponseWriter, statusCode int, data any) error {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)

	return json.NewEncoder(res).Encode(response{Data: data})
}
