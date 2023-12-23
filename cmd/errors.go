package main

import "net/http"

func (app *application) logError(err error) {
	app.logger.Error(err.Error())
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, errors map[string]string) {
	envelope := envelope{"error": errors}
	err := app.writeJSON(w, status, envelope)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, map[string]string{"server": message})
}
