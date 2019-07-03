package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (h *handler) index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	renderJSON(w, []byte(`
		{
			"module": "search",
			"version": "1.0.0", 
			"tagline": "You know, for search"
		}
	`), http.StatusOK)
}

func (h *notFoundHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	renderJSON(w, []byte(`
		{
			"message": "There's nothing here"
		}
	`), http.StatusNotFound)
}
