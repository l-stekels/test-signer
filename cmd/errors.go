package main

import (
	"fmt"
	"net/http"
)

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

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, map[string]string{"not_found": message})
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, map[string]string{"method_not_allowed": message})
}
