package data

import (
	"database/sql"
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
	query := `
	SELECT id, created_at, title, year, runtime, genres, version 
	FROM movies
	WHERE id = $1
	`
	var movie Movie
	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year,
		&movie.Runtime, pq.Array(&movie.Genres), &movie.Version,
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
	return nil
}

func (m MovieModel) Delete(id int64) error {
	return nil
}
