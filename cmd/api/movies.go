package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"vivekup3424/greenlight/internal/data"
)

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLogger.Println("Cant get ID(int)", err)
		http.Error(w, "Cant get ID", http.StatusBadRequest)
		return
	}
	movie := data.Movie{
		ID:      int64(numID),
		Title:   "Cocaina",
		Year:    2024,
		Runtime: 90,
		Genres:  []string{"drama", "more drama", "action"},
		Version: 1,
	}
	data, err := json.MarshalIndent(movie, "", "\t")
	if err != nil {
		app.errorLogger.Println("movie data marshalling:", err)
		http.Error(w, "Internal Server Error when getting the movie data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (app *application) createMoviesHandler(w http.ResponseWriter, r *http.Request) {
	//Declare an anonymous struct to hold the information that we
	//expect to be in the HTTP request body
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		app.errorLogger.Println("Decoding request body", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%+v\n", input)

	w.Write([]byte(`"message":"New movie created"`))
}
