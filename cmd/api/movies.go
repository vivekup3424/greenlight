package main

import (
	"fmt"
	"net/http"
)

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s := fmt.Sprintf("Showing details of movie : %s", id)
	w.Write([]byte(s))
}

func (app *application) createMoviesHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("POST request has been acknowledged"))
}
