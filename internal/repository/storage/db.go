package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DBStorage struct {
	db *sql.DB
}
