package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"test-signer.stekels.lv/internal/database/repositories"
	"test-signer.stekels.lv/internal/services"
	"testing"
)

func newTestApplication(t *testing.T, signatureRepo repositories.SignatureRepository) *application {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	return &application{
		logger:           logger,
		signatureService: services.NewSignatureService(logger, signatureRepo),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	t.Helper()
	ts := httptest.NewTLSServer(h)

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	res, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return res.StatusCode, res.Header, string(body)
}

func (ts *testServer) post(t *testing.T, urlPath string, payload any) (int, http.Header, string) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	jsonPayload = append(jsonPayload, '\n')
	res, err := ts.Client().Post(ts.URL+urlPath, "application/json", bytes.NewReader(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return res.StatusCode, res.Header, string(body)
}

type MockSignatureRepository struct {
	storage map[string]mockQuestions
}
type mockQuestions struct {
	Signature     string
	UserJWT       string
	QuestionsJson string
}

func (m MockSignatureRepository) Insert(signature string, userJWT string, questionsJson string) error {
	m.storage[signature] = mockQuestions{
		signature,
		userJWT,
		questionsJson,
	}

	return nil
}
