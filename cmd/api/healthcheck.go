package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	js, err := json.Marshal(data)
	if err != nil {
		app.errorLogger.Println("healthcheck data marshalling:", err)
		http.Error(w, "Internal Server Error when getting the healthcheck data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
