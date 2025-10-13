package config

import (
	"fmt"
	"net/url"
	"os"
)

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Config struct {
	DB *DB
}

func NewConfig() *Config {
	return &Config{
		DB: &DB{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
		},
	}
}

func (db *DB) GetDSN() string {
	if db == nil {
		return ""
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db.Host, db.Port, url.QueryEscape(db.User), url.QueryEscape(db.Password), db.DBName)
}
