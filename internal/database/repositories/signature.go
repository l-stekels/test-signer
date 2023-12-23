package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

var ErrDuplicateSignature = errors.New("duplicate signature")

type SignatureRepository interface {
	Insert(signature string, userJWT string, questions string) error
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
