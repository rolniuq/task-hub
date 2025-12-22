package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDB_GetDSN(t *testing.T) {
	db := &DB{
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
	}

	dsn := db.GetDSN()

	assert.Contains(t, dsn, "host=localhost")
	assert.Contains(t, dsn, "port=5432")
	assert.Contains(t, dsn, "user=testuser")
	assert.Contains(t, dsn, "password=testpass")
	assert.Contains(t, dsn, "dbname=testdb")
	assert.Contains(t, dsn, "sslmode=disable")
}

func TestDB_GetDSN_NilDB(t *testing.T) {
	var db *DB
	dsn := db.GetDSN()
	assert.Equal(t, "", dsn)
}

func TestDB_GetDSN_WithSpecialChars(t *testing.T) {
	db := &DB{
		Host:     "localhost",
		Port:     "5432",
		User:     "user@domain",
		Password: "pass#word!",
		DBName:   "testdb",
	}

	dsn := db.GetDSN()

	assert.NotEmpty(t, dsn)
	assert.Contains(t, dsn, "host=localhost")
}

func TestConfig_Fields(t *testing.T) {
	cfg := &Config{
		Port:         "8080",
		NatsUrl:      "nats://localhost:4222",
		JWTSecret:    "secret",
		DB:           &DB{Host: "localhost"},
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "nats://localhost:4222", cfg.NatsUrl)
	assert.Equal(t, "secret", cfg.JWTSecret)
	assert.NotNil(t, cfg.DB)
	assert.Equal(t, 15*time.Second, cfg.ReadTimeout)
	assert.Equal(t, 15*time.Second, cfg.WriteTimeout)
	assert.Equal(t, 60*time.Second, cfg.IdleTimeout)
}

func TestConfigModule(t *testing.T) {
	assert.NotNil(t, ConfigModule)
}
