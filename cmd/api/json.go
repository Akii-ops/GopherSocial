package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())

}

func wrtieJSON(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func NoContentResponse(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {

	// 限制读取的大小，不超过1M字节
	maxBytes := 1_048_578 // 1 M
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)

	// 禁止不明字段
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

type EnvelopeData struct {
	Data any `json:"data"`
}

type EnvelopeErr struct {
	Error string `json:"error"`
}

func wrtieJSONError(w http.ResponseWriter, status int, message string) error {

	return wrtieJSON(w, status, &EnvelopeErr{message})

}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {

	return wrtieJSON(w, status, &EnvelopeData{Data: data})
}
