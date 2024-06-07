package data

import "time"

type Movie struct {
	ID        int64     `json:"id"` //unique integer id for each movie
	CreatedAt time.Time `json:"-"`  //timestamp created for movie automatically when add to database
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   int32     `json:"runtime,omitempty"` //length(runtime) of movie
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}
