package data

import (
	"database/sql"
	"os"
)

// GetDB makes connection to database
func GetDB() (*sql.DB, error) {
	dbAddr := os.Getenv("DBADDRESS")

	db, err := sql.Open("mssql", dbAddr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
