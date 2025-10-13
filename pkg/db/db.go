package db

import (
	"database/sql"
	"taskhub/config"
)

type DB struct {
	conn *sql.DB
}

func NewDB(config *config.Config) *DB {
	conn, err := sql.Open("postgres", config.DB.GetDSN())
	if err != nil {
		return nil
	}

	return &DB{conn: conn}
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) GetConnection() *sql.DB {
	if db == nil {
		return nil
	}

	return db.conn
}
