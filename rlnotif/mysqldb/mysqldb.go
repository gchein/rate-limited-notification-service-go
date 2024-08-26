package mysqldb

import (
	"time"
)

type DB struct {
	// db *sql.DB

	Now func() time.Time
}
