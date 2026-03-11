package main

import (
	"log"
	"net/http"
)

// we differentiate the error packages/messages for client and developer

// for developer = use log
// for client = write into the response

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error: %s, path: %s, error: %s", r.Method, r.URL, err)

	writeJSONError(w, http.StatusInternalServerError, "Internal server error", err.Error())
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request error: %s, path: %s, error: %s", r.Method, r.URL, err)

	writeJSONError(w, http.StatusBadRequest, "Validation failed", err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not found error: %s, path: %s, error: %s", r.Method, r.URL, err)

	writeJSONError(w, http.StatusNotFound, "Resource not found", err.Error())
}
