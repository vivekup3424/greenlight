package data

import (
	"database/sql"
	"time"
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
	return nil
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

func (m MovieModel) Update(movie *Movie) error {
	return nil
}

func (m MovieModel) Delete(id int64) error {
	return nil
}
