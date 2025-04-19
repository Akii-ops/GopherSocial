package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	wrtieJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	wrtieJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponnse(w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	wrtieJSONError(w, http.StatusNotFound, "response not found")
}
