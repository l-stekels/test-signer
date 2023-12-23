package main

import (
	"net/http"
	"test-signer.stekels.lv/internal/transport"
	"test-signer.stekels.lv/internal/validator"
)

func (app *application) pingHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{
		"message": "Pong!",
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createSignatureHandler(w http.ResponseWriter, r *http.Request) {
	var input transport.CreateSignatureRequest
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	if v.Validate(input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, transport.NewCreateSignatureResponse("test_signature"))
}

func (app *application) verifySignatureHandler(w http.ResponseWriter, r *http.Request) {
	// so I need to read the request body first into a struct
	// {
	//	user_jwt: string,
	//	data: [
	//		{
	//			question: string,
	//			answer: string,
	//		},
	//		{}
	//	]
	//}
}
