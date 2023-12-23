package main

import "net/http"

type CreateSignatureRequest struct {
	UserJWT string     `json:"user_jwt"`
	Data    []Question `json:"data"`
}

type Question struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type CreateSignatureResponse struct {
	Signature string `json:"test_signature"`
}

func (app *application) pingHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{
		"message": "Pong!",
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createSignatureHandler(w http.ResponseWriter, r *http.Request) {
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
