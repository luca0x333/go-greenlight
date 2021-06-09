package main

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// readIDParam method retrives the "id" URL parameter from the current context request,
// then convert it to an interger and return it. If not successful returns an error.
func (app *application) readIDParam(r *http.Request) (int64, error) {
	// Any interpolated URL parameters will be stored in the request context r.Context().
	// Use ParamsFromContext() to retrieve a slice containing these parameters.
	params := httprouter.ParamsFromContext(r.Context())

	// Use ByName() method to retrieve the value of "id" parameter from the slice.
	// The value returned by ByName() is a string so we need to convert it to a base 10 integer.
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// writeJSON method return a json response.
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	// json.Marshal() returns a []byte slice containing the encoded JSON.
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	// We loop through the header map and add each header to the http.ResponseWriter header map.
	// Go doesn't throw any error if we range over a nil map.
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
