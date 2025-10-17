package config

import (
	"fmt"
	"net/url"
	"os"
	"time"
)

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Config struct {
	Port         string
	NatsUrl      string
	JWTSecret    string
	DB           *DB
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func NewConfig() *Config {
	return &Config{
		Port:      os.Getenv("PORT"),
		NatsUrl:   os.Getenv("NATS_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		DB: &DB{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
		},
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func (db *DB) GetDSN() string {
	if db == nil {
		return ""
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db.Host, db.Port, url.QueryEscape(db.User), url.QueryEscape(db.Password), db.DBName)
}
