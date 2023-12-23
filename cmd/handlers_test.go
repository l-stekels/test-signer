package main

import (
	"net/http"
	"testing"
)

func TestPingGet(t *testing.T) {
	app := newTestApplication(t)

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
	app := newTestApplication(t)
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

//func TestCreateSignature(t *testing.T) {
//	app := newTestApplication(t)
//
//	ts := newTestServer(t, app.routes())
//	defer ts.Close()
//
//	code, _, body := ts.post(t, "/create-signature")
//
//	if code != http.StatusOK {
//		t.Errorf("want %d; got %d", http.StatusOK, code)
//	}
//
//	if string(body) != `{"message":"Pong!"}` {
//		t.Errorf("want body to equal %q", `{"message":"Pong!"}`)
//	}
//}
