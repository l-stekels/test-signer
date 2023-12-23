package main

import (
	"database/sql"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
	"os"
	"test-signer.stekels.lv/internal/database/repositories"
	"test-signer.stekels.lv/internal/services"
	"time"
)

type application struct {
	config           config
	logger           *slog.Logger
	signatureService *services.SignatureService
}

type config struct {
	env         string
	db          dbConfig
	redisConfig redisConfig
}

type dbConfig struct {
	username string
	password string
	database string
	host     string
}
type redisConfig struct {
	host string
	port string
}

func main() {
	cfg := config{
		env: os.Getenv("ENV"),
		db: dbConfig{
			username: os.Getenv("MYSQL_USER"),
			password: os.Getenv("MYSQL_PASSWORD"),
			database: os.Getenv("MYSQL_DATABASE"),
			host:     os.Getenv("MYSQL_HOST"),
		},
		redisConfig: redisConfig{
			host: os.Getenv("REDIS_HOST"),
			port: os.Getenv("REDIS_PORT"),
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    fmt.Sprintf("%s:%s", cfg.redisConfig.host, cfg.redisConfig.port),
	})
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}(client)
	locker := redislock.New(client)

	db, err := openDB(cfg)
	if err != nil {
		os.Exit(1)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}(db)

	signatureRepo := repositories.NewMySQLSignatureRepository(db)
	redisLocker := services.NewRedisLocker(locker)
	app := &application{
		config:           cfg,
		logger:           logger,
		signatureService: services.NewSignatureService(logger, signatureRepo, redisLocker),
	}
	server := &http.Server{
		Addr:         ":4000",
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Starting server on port 400", slog.String("env", cfg.env))
	err = server.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.User = cfg.db.username
	mysqlConfig.Passwd = cfg.db.password
	mysqlConfig.DBName = cfg.db.database
	mysqlConfig.Net = "tcp"
	mysqlConfig.Addr = cfg.db.host
	mysqlConfig.ParseTime = true

	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(15 * time.Minute)
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
