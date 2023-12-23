package main

import (
	"errors"
	"net/http"
	"test-signer.stekels.lv/internal/models"
	"test-signer.stekels.lv/internal/services"
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
	if v.ValidateCreateSignatureRequest(input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	signatureModel := &models.Signature{
		UserJWT: input.UserJWT,
	}
	for _, question := range input.Questions {
		signatureModel.Questions = append(signatureModel.Questions, models.Question{
			Body:   question.Body,
			Answer: question.Answer,
		})
	}

	signature, err := app.signatureService.Create(*signatureModel)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	err = app.writeJSON(w, http.StatusCreated, transport.NewCreateSignatureResponse(signature))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) verifySignatureHandler(w http.ResponseWriter, r *http.Request) {
	var input transport.VerifySignatureRequest
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	if v.ValidateVerifySignatureRequest(input); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	signatureModel, err := app.signatureService.Get(input.Signature)
	if errors.Is(err, services.ErrSignatureNotFound) {
		app.notFoundResponse(w, r)

		return
	}
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}
	if input.UserJWT != signatureModel.UserJWT {
		app.notFoundResponse(w, r)

		return
	}

	var questions []transport.Question
	for _, question := range signatureModel.Questions {
		questions = append(questions, transport.Question{
			Body:   question.Body,
			Answer: question.Answer,
		})
	}

	err = app.writeJSON(w, http.StatusOK, transport.NewVerifySignatureResponse(questions, signatureModel.CreatedAt))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
