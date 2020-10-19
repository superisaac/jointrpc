package jsonrpc

import (
	"log"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, r *http.Request, err error, status int, message string) {
	log.Printf("HTTP error: %s %d", err.Error(), status)
	w.WriteHeader(status)
	w.Write([]byte(message))
}

