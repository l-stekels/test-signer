package services

import (
	"database/sql"
	"log/slog"
	"test-signer.stekels.lv/internal/transport"
)

type SignatureService struct {
	logger *slog.Logger
	db     *sql.DB
}

func NewSignatureService(logger *slog.Logger, db *sql.DB) *SignatureService {
	return &SignatureService{
		logger: logger,
		db:     db,
	}
}

func (s *SignatureService) Create(input transport.CreateSignatureRequest) (string, error) {
	// 1. create signature
	// 2. create questions
	// 3. create answers

	// store all that in the dataabase

	return "test_signature", nil
}
