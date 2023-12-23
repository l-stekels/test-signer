package services

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/hashstructure/v2"
	"log/slog"
	"strconv"
	"test-signer.stekels.lv/internal/database/repositories"
	"test-signer.stekels.lv/internal/models"
)

type SignatureService struct {
	logger        *slog.Logger
	signatureRepo repositories.SignatureRepository
}

func NewSignatureService(logger *slog.Logger, signatureRepository repositories.SignatureRepository) *SignatureService {
	return &SignatureService{
		logger:        logger,
		signatureRepo: signatureRepository,
	}
}

func (s *SignatureService) Create(model models.Signature) (string, error) {
	// TODO: Create a lock here and unlock only after everything is OK
	signature, err := createSignature(model)
	if err != nil {
		return "", err
	}
	questionsJson, err := json.Marshal(model.Questions)
	err = s.signatureRepo.Insert(signature, model.UserJWT, string(questionsJson))
	if errors.Is(err, repositories.ErrDuplicateSignature) {
		s.logger.Info("same signature detected", "signature", signature)

		return signature, nil
	}
	if err != nil {
		return "", err
	}

	return signature, nil
}

func createSignature(model models.Signature) (string, error) {
	hash, err := hashstructure.Hash(model, hashstructure.FormatV2, nil)
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(hash, 10), nil
}
