package models

import (
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Signature struct {
	ID        int `hash:"ignore"`
	UserJWT   string
	Signature string    `hash:"ignore"`
	CreatedAt time.Time `hash:"ignore"`
	Questions []Question
}

type Question struct {
	Body   string
	Answer string
}
