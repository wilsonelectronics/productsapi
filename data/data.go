package data

import (
	"database/sql"
	"os"

	// Justification: No main.go file, and this is the only file that 'technically' uses this package
	_ "github.com/denisenkom/go-mssqldb"
)

// GetDB makes connection to database
func GetDB() (db *sql.DB, err error) {
	if db, err = sql.Open("mssql", os.Getenv("DBADDRESS")); err != nil {
		return nil, err
	}
	return db, nil
}
