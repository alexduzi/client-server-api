package database

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	TIMEOUT_EXCHANGE_DB time.Duration = time.Millisecond * 10 // 10 ms
)

var (
	db   *sql.DB
	stmt *sql.Stmt
)

func createDb() {
	file, err := os.OpenFile("./cotacao.db", os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("failed to open/create file: %v", err)
	}
	defer file.Close()
}

func InitializeDb() {
	createDb()

	db, err := sql.Open("sqlite3", "./cotacao.db?cache=shared&mode=rwc&_journal_mode=WAL&_synchronous=NORMAL&_timeout=5000")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS exchange (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		url VARCHAR(100) NOT NULL,
		data TEXT NOT NULL,
		createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("error creating table: %q: %s\n", err, sqlStmt)
	}

	sqlInsert := `INSERT INTO exchange (url, data) VALUES(?, ?)`
	stmt, err = db.Prepare(sqlInsert)
	if err != nil {
		log.Fatalf("error preparing statement: %v", err)
	}
}

func InsertExchange(ctx context.Context, url, json string) error {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT_EXCHANGE_DB)
	defer cancel()

	_, err := stmt.ExecContext(ctx, url, json)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
		return ctx.Err()
	default:
		log.Println("exchange inserted!")
	}

	return nil
}

func CloseDb() {
	if stmt != nil {
		stmt.Close()
	}
	if db != nil {
		db.Close()
	}
}
