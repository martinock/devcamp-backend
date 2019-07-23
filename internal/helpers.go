package internal

import (
	"net/http"
)

func renderJSON(w http.ResponseWriter, data []byte, status int) {

	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Most of the API HTTP response is a JSON format
	w.Header().Set("Content-Type", "application/json")

	// HTTP status (200 OK, 404 Not Found, 500 Internal Server Error, etc.)
	w.WriteHeader(status)

	// The actual data
	w.Write(data)
}
