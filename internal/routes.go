package main

import "net/http"

func (app *application) routes() *http.ServeMux{
    router := http.NewServeMux()
    return router
}

