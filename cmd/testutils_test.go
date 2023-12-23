package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"test-signer.stekels.lv/internal/database/repositories"
	"test-signer.stekels.lv/internal/models"
	"test-signer.stekels.lv/internal/services"
	"testing"
)

func newTestApplication(t *testing.T, signatureRepo repositories.SignatureRepository) *application {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	return &application{
		logger:           logger,
		signatureService: services.NewSignatureService(logger, signatureRepo, &MockLocker{}),
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
	storage map[string]models.Signature
}

func (m *MockSignatureRepository) Insert(signature string, userJWT string, questionsJson string) error {
	var questions []models.Question
	err := json.Unmarshal([]byte(questionsJson), &questions)
	if err != nil {
		return err
	}

	m.storage[signature] = models.Signature{
		Signature: signature,
		UserJWT:   userJWT,
		Questions: questions,
	}

	return nil
}

func (m *MockSignatureRepository) GetBySignature(signature string) (*models.Signature, error) {
	model, ok := m.storage[signature]
	if !ok {
		return nil, repositories.ErrRecordNotFound
	}

	return &model, nil
}

type MockLocker struct{}

func (m *MockLocker) GetLock(key string) error {
	return nil
}
func (m *MockLocker) ReleaseLock(key string) error {
	return nil
}
