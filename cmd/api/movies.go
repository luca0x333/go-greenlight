package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// createMovieHandler "POST /v1/movies"
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

// showMovieHandler "GET /v1/movies/:id"
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Any interpolated URL parameters will be stored in the request context.
	// Use ParamsFromContext() to retrieve a slice containing these parameters.
	params := httprouter.ParamsFromContext(r.Context())

	// Use ByName() method to retrieve the value of "id" parameter from the slice.
	// The value returned by ByName() is a string so we need to convert it to a base 10 integer.
	id, err := strconv.ParseInt(params.ByName("id"), 10, 60)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "show the details of movie %d\n", id)
}
