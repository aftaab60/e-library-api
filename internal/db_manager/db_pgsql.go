package db_manager

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

// wrapper to customise db query methods below.
// If we need to support different database types, change this to interface and make dependency injection as needed for different db types
type DB struct {
	db *sql.DB
}

var once sync.Once
var db *DB

func InitPgsqlConnection() *DB {
	once.Do(func() {
		//TODO: setup to read from config file. For example config.yml
		//for now, hardcoded it with values in docker container
		dbHost := "localhost"
		dbPort := "5432"
		dbUser := "userdev"
		dbPassword := "dev123"
		dbName := "db_pgsql"

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName)

		var err error
		dbpgsql, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("Error opening database: %v", err)
		}

		if err = dbpgsql.Ping(); err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}

		log.Println("Successfully connected to PostgreSQL database")
		db = &DB{db: dbpgsql}
	})
	return db
}

func (d *DB) CreateRecord(ctx context.Context, query string, args ...interface{}) *sql.Row {
	tx := GetTransactionFromContext(ctx)
	if tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *DB) GetRecord(ctx context.Context, query string, args ...interface{}) *sql.Row {
	tx := GetTransactionFromContext(ctx)
	if tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *DB) UpdateRecord(ctx context.Context, query string, args ...interface{}) *sql.Row {
	tx := GetTransactionFromContext(ctx)
	if tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *DB) DeleteRecord(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	tx := GetTransactionFromContext(ctx)
	if tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}
	return d.db.ExecContext(ctx, query, args...)
}

func (d *DB) Begin() (*sql.Tx, error) {
	return d.db.Begin()
}

func CloseDB() {
	if db != nil {
		err := db.db.Close()
		if err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}
}

// methods for transaction
type ItxDB interface {
	Begin() (*sql.Tx, error)
}
