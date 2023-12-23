package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/ping", app.pingHandler)
	router.HandlerFunc(http.MethodPost, "/signature", app.createSignatureHandler)
	router.HandlerFunc(http.MethodPost, "/signature/verify", app.verifySignatureHandler)

	return router
}
