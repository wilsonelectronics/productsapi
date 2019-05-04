package data

import (
	"database/sql"
	"os"

	// Justification: No main.go file, and this is the only file that 'technically' uses this package
	_ "github.com/denisenkom/go-mssqldb"
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
