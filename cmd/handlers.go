package main

import "net/http"

func (app *application) pingHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{
		"message": "Pong!",
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
