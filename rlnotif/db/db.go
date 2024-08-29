package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
	Ping() error
	Begin() (*sql.Tx, error)
}

func NewMySQLStorage(cfg *mysql.Config) (*DB, error) {
	var err error
	var newDB DB
	newDB, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	return &newDB, nil
}
