package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

// DB interface should receive *sql.DB or *sql.Tx
type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
}

func NewMySQLStorage(cfg *mysql.Config) (*sql.DB, error) {
	DB, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	return DB, nil
}
