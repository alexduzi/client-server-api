package database

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func createDb() {
	file, err := os.OpenFile("./exchange.db", os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("failed to open/create file: %v", err)
	}
	defer file.Close()
}

func openConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./exchange.db")
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

func InitializeDb() {
	createDb()

	db, _ := openConnection()
	defer db.Close()

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS exchange (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		url VARCHAR(100) NOT NULL,
		data TEXT NOT NULL,
		createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("error creating table: %q: %s\n", err, sqlStmt)
	}
}

func InsertExchange(ctx context.Context, url, json string) error {
	db, _ := openConnection()
	defer db.Close()

	sqlSmt := `INSERT INTO exchange (url, data) VALUES(?, ?)`

	insertPrepare, err := db.Prepare(sqlSmt)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	smt := tx.StmtContext(ctx, insertPrepare)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer smt.Close()

	_, err = smt.ExecContext(ctx, url, json)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}
