package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bsm/redislock"
	"github.com/mitchellh/hashstructure/v2"
	"log/slog"
	"strconv"
	"test-signer.stekels.lv/internal/database/repositories"
	"test-signer.stekels.lv/internal/models"
	"time"
)

var ErrSignatureNotFound = errors.New("signature not found")

type SignatureService struct {
	logger        *slog.Logger
	signatureRepo repositories.SignatureRepository
	locker        Locker
}

func NewSignatureService(logger *slog.Logger, signatureRepository repositories.SignatureRepository, locker Locker) *SignatureService {
	return &SignatureService{
		logger:        logger,
		signatureRepo: signatureRepository,
		locker:        locker,
	}
}

type Locker interface {
	GetLock(key string) error
	ReleaseLock(key string) error
}
type RedisLocker struct {
	locks  map[string]*redislock.Lock
	client *redislock.Client
}

func NewRedisLocker(client *redislock.Client) *RedisLocker {
	return &RedisLocker{
		locks:  make(map[string]*redislock.Lock),
		client: client,
	}
}
func (r *RedisLocker) GetLock(key string) error {
	ctx := context.Background()
	lock, err := r.client.Obtain(ctx, key, 10*time.Second, nil)
	if errors.Is(err, redislock.ErrNotObtained) {
		return err
	}
	if err != nil {
		return err
	}
	r.locks[key] = lock

	return nil
}
func (r *RedisLocker) ReleaseLock(key string) error {
	ctx := context.Background()
	lock, ok := r.locks[key]
	if !ok {
		return nil
	}
	err := lock.Release(ctx)
	if err != nil {
		return err
	}
	delete(r.locks, key)

	return nil
}

func (s *SignatureService) Create(model models.Signature) (string, error) {
	signature, err := createSignature(model)
	if err != nil {
		return "", err
	}
	err = s.locker.GetLock(signature)
	if err != nil {
		return "", err
	}
	defer s.locker.ReleaseLock(signature)

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

func (s *SignatureService) Get(signature string) (models.Signature, error) {
	model, err := s.signatureRepo.GetBySignature(signature)
	if errors.Is(err, repositories.ErrRecordNotFound) {
		return models.Signature{}, ErrSignatureNotFound
	}

	return *model, err
}

func createSignature(model models.Signature) (string, error) {
	hash, err := hashstructure.Hash(model, hashstructure.FormatV2, nil)
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(hash, 10), nil
}
