package main

import (
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
	id, err := strconv.ParseInt(params.ByName("id"), 64, 10)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}
