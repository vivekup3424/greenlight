package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

// Define a config struct to hold all the configuration settings for our application.
// For now, the only configuration settings will be the network port that we want the
// server to listen on, and the name of the current operating environment for the
// application (development, staging, production, etc.). We will read in these
// configuration settings from command-line flags when the application starts.
type config struct {
	port int
	env  string
}
// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as our build progresses.
type application struct {
	config config
	logger *log.Logger
}


func main() {
	//declate an instance of config struct
	var cfg config

	//read the value of port and env-commandline flags into the
	//config struct.
	//defaults are 4000 port and "development" environment
	flag.IntVar(&cfg.port, "port", 4000, "API Server Port")
	flag.StringVar(&cfg.env, "env", "development",
		"Environment (development|staging|production)")

	//logger to write message to stdout
	logger := log.New(os.Stdout, "INFO", log.Ldate|log.Ltime|log.Lshortfile)

	//an instance of the application struct
	app := application{
		config: cfg,
		logger: logger,
	}
	//declare a new router and a /v1/healthcheck route
	router := http.NewServeMux()
	router.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
    router.HandleFunc("POST /v1/movies", app.createMoviesHandler)
    router.HandleFunc("GET /v1/movies/{id}", app.showMovieHandler)

	//declare a http with some good timeout settings. which listens
	//on the provided with port, and the above router as the handler

	srv := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", cfg.port),
		Handler:      router,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	//Start the http server
	logger.Printf("starting the %s server on : %s", cfg.env, srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}

}
