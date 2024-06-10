package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

type Movie struct {
	ID        int64     `json:"id"` //unique integer id for each movie
	CreatedAt time.Time `json:"-"`  //timestamp created for movie automatically when add to database
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   int32     `json:"runtime,omitempty"` //length(runtime) of movie
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}
type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) GetALl() ([]*Movie, error) {
	query := `
	SELECT id, created_at,title,year,runtime,genres,version
	FROM movies
	ORDER BY id
	`
	//create an context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		log.Println("Error getting movies", err)
		return nil, err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()

	movies := []*Movie{}

	for rows.Next() {
		var movie Movie

		err = rows.Scan(
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK, then return the slice of movies.
	return movies, nil
}
func (m MovieModel) Insert(movie *Movie) error {
	query := `
	INSERT INTO movies (title,year,runtime,genres)
	VALUES ($1,$2,$3,$4)
	RETURNING id,created_at,version
	`
	//data := []interface{}{movie.Title, movie.Year, movie.Runtime,
	//pq.A}
	err := m.DB.QueryRow(query, movie.Title, movie.Year, movie.Runtime,
		pq.Array(movie.Genres)).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
	if err != nil {
		log.Println("creating movie inside greenlight database", err)
	} else {
		log.Printf("movie with id: %d created successfully inside greenlight database\n", movie.ID)
	}
	return err
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	// Use the context.WithTimeout() function to create a context.Context which carries a
	// 1-second timeout deadline. Note that we're using the empty context.Background()
	// as the 'parent' context.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// Importantly, use defer to make sure that we cancel the context before the Get()
	// method returns.
	defer cancel()
	query := `
	SELECT id, created_at, title, year, runtime, genres, version 
	FROM movies
	WHERE id = $1
	`
	var movie Movie
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Data row not found in database", err)
			return nil, ErrRecordNotFound
		} else {
			fmt.Println("Some unknown error occured", err)
			return nil, err
		}
	}
	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
	query := `
	UPDATE movies
	SET title = $1, year = $2, runtime = $3, genres = $4,version=version+1
	WHERE id = $5 and version=$6
	RETURNING version
	`
	err := m.DB.QueryRow(query, movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.ID, movie.Version).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			log.Println("Edit Conflict(version)", err)
			return ErrEditConflict
		default:
			log.Println("Updating movie", err)
			return err
		}
	} else {
		log.Println("Movie updated successfully")
	}
	return nil //no error
}

func (m MovieModel) Delete(id int64) error {
	query := `
	DELETE FROM movies
	WHERE id = $1
	`
	results, err := m.DB.Exec(query, id)
	if err != nil {
		fmt.Println("delete operation", err)
		return err
	}
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		fmt.Println("dont know about this fucking error", err)
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
