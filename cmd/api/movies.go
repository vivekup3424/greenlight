package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	movie, err := app.models.Movies.Get(int64(numID))
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.errorLogger.Println("getting movie", err)
			http.Error(w, "movie not found", http.StatusNotFound)
		} else {
			app.errorLogger.Println("unkown error getting movie", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	//movie := data.Movie{
	//	ID:      int64(numID),
	//	Title:   "Cocaina",
	//	Year:    2024,
	//	Runtime: 90,
	//	Genres:  []string{"drama", "more drama", "action"},
	//	Version: 1,
	//}
	data, err := json.Marshal(movie)
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
	//feeding the data on the database
	newMovie := data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}
	if err := app.models.Movies.Insert(&newMovie); err != nil {
		app.errorLogger.Println("Inserting movie into database", err)
		http.Error(w, "Database Insertion Errror", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(newMovie)
	if err != nil {
		app.errorLogger.Println("Converting the movie struct to json", err)
		w.Write([]byte(`"message":"movie marshalling to json failed`))
	}
	w.Write(js)
	w.Write([]byte(`"message":"New movie created"`))
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLogger.Println("Cant get ID(int)", err)
		http.Error(w, "Cant get ID", http.StatusBadRequest)
		return
	}
	movie, err := app.models.Movies.Get(int64(numID))
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.errorLogger.Println("getting movie", err)
			http.Error(w, "movie not found", http.StatusNotFound)
		} else {
			app.errorLogger.Println("unkown error getting movie", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	//get the new value for movie update
	var input struct {
		Title   *string  `json:"title"`
		Year    *int32   `json:"year"`
		Runtime *int32   `json:"runtime"`
		Genres  []string `json:"genres"`
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&input) //some golang nuances
	if err != nil {
		app.errorLogger.Println("Decoding request body", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	//copy the values from input to the movie pointer
	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}
	err = app.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.errorLogger.Println("edit conflict", err)
			http.Error(w, "unable to update the record due to edit conflict, please try again", http.StatusConflict)
			return
		default:
			app.errorLogger.Println("updating movie id=", movie.ID, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("movie updated successfully")) //this also returns an errr, but I dont know what to do with that error
	fmt.Fprintf(w, "%+v", movie)
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	numID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLogger.Println("Cant get ID(int)", err)
		http.Error(w, "Cant get ID", http.StatusBadRequest)
		return
	}
	err = app.models.Movies.Delete(int64(numID))
	if err == data.ErrRecordNotFound {
		app.errorLogger.Println("movie id not found", err)
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	} else if err != nil {
		app.errorLogger.Println("failed delete operation", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	app.infoLogger.Printf("movie with id: %v deleted successfully from database\n", numID)
	w.Write([]byte("movie deleted successfully"))
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string
		Genres  []string
		Filters data.Filters
	}
	//parse and get the params from the query string
	queryString := r.URL.Query()

	var title string
	defaultTitle := ""
	//get the title from querystring
	titleFromQuery := queryString.Get("title")
	if titleFromQuery == "" {
		title = defaultTitle
	} else {
		title = titleFromQuery
	}
	input.Title = title

	//get the genres from the querystring
	var genres []string
	defaultGenres := []string{}
	genresFromQuery := queryString.Get("genres")
	if genresFromQuery == "" {
		genres = defaultGenres
	} else {
		genres = strings.Split(genresFromQuery, ",")
	}
	input.Genres = genres

	//get the page number
	input.Filters.Page = app.readInt(queryString, "page", 1)
	input.Filters.PageSize = app.readInt(queryString, "page_size", 20)

	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply a ascending sort on movie ID).
	input.Filters.Sort = app.readString(queryString, "sort", "id")

	movies, err := app.models.Movies.GetALl()
	if err != nil {
		app.errorLogger.Println("getting movies", err)
		http.Error(w, "error when getting movies", http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(movies)
	if err != nil {
		app.errorLogger.Println("marshalln=ing movie, converting to json", err)
		http.Error(w, "error when getting movies", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application-json")
	w.Write(js)
	//dump on the content on the response writer
	fmt.Fprintf(w, "%+v\n", input)
}
