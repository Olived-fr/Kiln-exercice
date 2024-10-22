package api

import (
	"net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

// Handle wraps a http.HandlerFunc with error handling.
func Handle(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			writeError(w, r, err)
		}
	}
}
