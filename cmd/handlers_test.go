package main

import (
	"net/http"
	"test-signer.stekels.lv/internal/transport"
	"testing"
)

func TestPingGet(t *testing.T) {
	app := newTestApplication(t, nil)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	expectedBody := `{"message":"Pong!"}`
	if string(body) != expectedBody {
		t.Errorf("want body to equal %q, got %q", expectedBody, body)
	}
}

func TestPingPost(t *testing.T) {
	app := newTestApplication(t, nil)
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	payload := struct {
		Message string `json:"message"`
	}{
		Message: "Pong!",
	}

	code, _, body := ts.post(t, "/ping", payload)

	if code != http.StatusMethodNotAllowed {
		t.Errorf("want %d; got %d", http.StatusMethodNotAllowed, code)
	}
	expectedBody := `{"error":{"method_not_allowed":"the POST method is not supported for this resource"}}`
	if string(body) != expectedBody {
		t.Errorf("want body to equal %q, got %q", expectedBody, body)
	}
}

func TestCreateSignatureErrors(t *testing.T) {
	app := newTestApplication(t, nil)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name               string
		payload            interface{}
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "Empty payload",
			payload:            struct{}{},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedBody:       `{"error":{"questions":"must be provided","user_jwt":"must be provided"}}`,
		},
		{
			name: "Empty questions",
			payload: struct {
				UserJWT string `json:"user_jwt"`
			}{
				"this_is_a_jwt",
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedBody:       `{"error":{"questions":"must be provided"}}`,
		},
		{
			name: "Questions without answers",
			payload: struct {
				UserJWT   string               `json:"user_jwt"`
				Questions []transport.Question `json:"questions"`
			}{
				"this_is_a_jwt",
				[]transport.Question{
					{
						Body: "What is your name?",
					},
					{
						Body: "How are you?",
					},
				},
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedBody:       `{"error":{"question[0]":"question answer must be provided","question[1]":"question answer must be provided"}}`,
		},
		{
			name: "Answers without questions",
			payload: struct {
				UserJWT   string               `json:"user_jwt"`
				Questions []transport.Question `json:"questions"`
			}{
				"this_is_a_jwt",
				[]transport.Question{
					{
						Answer: "42",
					},
					{
						Answer: "A",
					},
				},
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedBody:       `{"error":{"question[0]":"question body must be provided","question[1]":"question body must be provided"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.post(t, "/signature", tt.payload)
			if code != tt.expectedStatusCode {
				t.Errorf("want %d; got %d", tt.expectedStatusCode, code)
			}
			if body != tt.expectedBody {
				t.Errorf("want body to equal %q, got %q", tt.expectedBody, body)
			}
		})
	}
}

func TestCreateSignature(t *testing.T) {
	app := newTestApplication(t, MockSignatureRepository{make(map[string]mockQuestions)})
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	payload := struct {
		UserJWT   string               `json:"user_jwt"`
		Questions []transport.Question `json:"questions"`
	}{
		"this_is_a_jwt",
		[]transport.Question{
			{
				Body:   "What is your name?",
				Answer: "John Doe",
			},
			{
				Body:   "The answer to life, the universe, and everything?",
				Answer: "42",
			},
		},
	}

	code, _, body := ts.post(t, "/signature", payload)

	if code != http.StatusCreated {
		t.Errorf("want %d; got %d", http.StatusCreated, code)
	}
	expectedBody := `{"data":{"signature":"14575973885829822288"}}`
	if body != expectedBody {
		t.Errorf("want body to equal %q, got %q", expectedBody, body)
	}
}
