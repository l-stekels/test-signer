package main

import (
	"net/http"
	"test-signer.stekels.lv/internal/models"
	"test-signer.stekels.lv/internal/transport"
	"testing"
	"time"
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
	app := newTestApplication(t, &MockSignatureRepository{make(map[string]models.Signature)})
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

func TestValidateSignatureErrors(t *testing.T) {
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
			expectedBody:       `{"error":{"signature":"must be provided","user_jwt":"must be provided"}}`,
		},
		{
			name: "Empty signature",
			payload: struct {
				UserJWT string `json:"user_jwt"`
			}{
				"this_is_a_jwt",
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedBody:       `{"error":{"signature":"must be provided"}}`,
		},
		{
			name: "Empty signature",
			payload: struct {
				Signature string `json:"signature"`
			}{
				"123456",
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedBody:       `{"error":{"user_jwt":"must be provided"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.post(t, "/signature/verify", tt.payload)
			if code != tt.expectedStatusCode {
				t.Errorf("want %d; got %d", tt.expectedStatusCode, code)
			}
			if body != tt.expectedBody {
				t.Errorf("want body to equal %q, got %q", tt.expectedBody, body)
			}
		})
	}
}

func TestVerifySignatureExists(t *testing.T) {
	repoStorage := make(map[string]models.Signature)
	repoStorage["14575973885829822288"] = models.Signature{
		Signature: "14575973885829822288",
		UserJWT:   "this_is_a_jwt",
		Questions: []models.Question{
			{
				Body:   "What is the answer to the universe and everything?",
				Answer: "42",
			},
		},
		CreatedAt: time.Date(2023, 12, 23, 16, 00, 00, 1, time.UTC),
	}
	app := newTestApplication(t, &MockSignatureRepository{repoStorage})
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	payload := struct {
		UserJWT   string `json:"user_jwt"`
		Signature string `json:"signature"`
	}{
		"this_is_a_jwt",
		"14575973885829822288",
	}

	code, _, body := ts.post(t, "/signature/verify", payload)

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}
	expectedBody := `{"data":{"answers":[{"question":"What is the answer to the universe and everything?","answer":"42"}],"timestamp":"2023-12-23T16:00:00.000000001Z"}}`
	if body != expectedBody {
		t.Errorf("want body to equal %q, got %q", expectedBody, body)
	}
}

func TestVerifySignatureExistsDifferentJWT(t *testing.T) {
	repoStorage := make(map[string]models.Signature)
	repoStorage["14575973885829822288"] = models.Signature{
		Signature: "14575973885829822288",
		UserJWT:   "different_user_jwt",
		Questions: []models.Question{
			{
				Body:   "What is the answer to the universe and everything?",
				Answer: "42",
			},
		},
		CreatedAt: time.Date(2023, 12, 23, 16, 00, 00, 1, time.UTC),
	}
	app := newTestApplication(t, &MockSignatureRepository{repoStorage})
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	payload := struct {
		UserJWT   string `json:"user_jwt"`
		Signature string `json:"signature"`
	}{
		"this_is_a_jwt",
		"14575973885829822288",
	}

	code, _, body := ts.post(t, "/signature/verify", payload)

	if code != http.StatusNotFound {
		t.Errorf("want %d; got %d", http.StatusNotFound, code)
	}
	expectedBody := `{"error":{"not_found":"the requested resource could not be found"}}`
	if body != expectedBody {
		t.Errorf("want body to equal %q, got %q", expectedBody, body)
	}
}

func TestVerifySignatureNotFound(t *testing.T) {
	repoStorage := make(map[string]models.Signature)
	app := newTestApplication(t, &MockSignatureRepository{repoStorage})
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	payload := struct {
		UserJWT   string `json:"user_jwt"`
		Signature string `json:"signature"`
	}{
		"this_is_a_jwt",
		"14575973885829822288",
	}

	code, _, body := ts.post(t, "/signature/verify", payload)

	if code != http.StatusNotFound {
		t.Errorf("want %d; got %d", http.StatusNotFound, code)
	}
	expectedBody := `{"error":{"not_found":"the requested resource could not be found"}}`
	if body != expectedBody {
		t.Errorf("want body to equal %q, got %q", expectedBody, body)
	}
}
