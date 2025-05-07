package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorw("internal error", "method", r.Method, "paht", r.URL.Path, "error", err.Error())
	wrtieJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Warnw("forbidden error", "method", r.Method, "paht", r.URL.Path, "error", err.Error())
	wrtieJSONError(w, http.StatusForbidden, "forbidden")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Warnf("bad request", "method", r.Method, "paht", r.URL.Path, "error", err.Error())
	wrtieJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponnse(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Warnf("not found ", "method", r.Method, "paht", r.URL.Path, "error", err.Error())
	wrtieJSONError(w, http.StatusNotFound, "response not found")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorf("conflict respons", "method", r.Method, "paht", r.URL.Path, "error", err.Error())
	wrtieJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Warnf("unauthorized error", "method", r.Method, "paht", r.URL.Path, "error", err.Error())
	wrtieJSONError(w, http.StatusUnauthorized, "unauthorized")
}

// not used
// func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

// 	app.logger.Warnf("unauthorized basci error", "method", r.Method, "paht", r.URL.Path, "error", err.Error())

// 	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
// 	w.WriteHeader(http.StatusUnauthorized)
// 	wrtieJSONError(w, http.StatusConflict, "unauthorized")
// }

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryafter string) error {
	app.logger.Warnf("rate limit exceeded", "method", r.Method, "path", r.URL.Path)
	w.Header().Set("Retry-After", retryafter)
	return wrtieJSON(w, http.StatusTooManyRequests, "rate limit exceeded, retry after :"+retryafter)
}
