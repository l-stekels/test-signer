package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-sql-driver/mysql"
	"strings"
	"test-signer.stekels.lv/internal/models"
	"time"
)

var ErrDuplicateSignature = errors.New("duplicate signature")
var ErrRecordNotFound = errors.New("record not found")

type SignatureRepository interface {
	Insert(signature string, userJWT string, questions string) error
	GetBySignature(signature string) (*models.Signature, error)
}

type MySQLSignatureRepository struct {
	Conn *sql.DB
}

func NewMySQLSignatureRepository(conn *sql.DB) *MySQLSignatureRepository {
	return &MySQLSignatureRepository{Conn: conn}
}

func (r *MySQLSignatureRepository) Insert(signature string, userJWT string, questions string) error {
	query := `INSERT INTO signatures (signature, user_jwt, questions) VALUES (?, ?, ?)`
	args := []interface{}{signature, userJWT, questions}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := r.Conn.ExecContext(ctx, query, args...)
	if err != nil {
		var mySQLError *mysql.MySQLError
		// Check for mysql error 1062 (duplicate entry for email)
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "signature_unique") {
				return ErrDuplicateSignature
			}
		}

		return err
	}
	return nil
}

func (r *MySQLSignatureRepository) GetBySignature(signature string) (*models.Signature, error) {
	query := `SELECT s.user_jwt, s.created_at, s.questions FROM signatures s WHERE signature = ?`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var signatureModel models.Signature
	var questionsJson string
	err := r.Conn.QueryRowContext(ctx, query, signature).Scan(
		&signatureModel.UserJWT,
		&signatureModel.CreatedAt,
		&questionsJson,
	)
	var questions []models.Question
	err = json.Unmarshal([]byte(questionsJson), &questions)
	if err != nil {
		return nil, err
	}
	signatureModel.Questions = questions
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &signatureModel, nil
}
