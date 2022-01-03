package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// connect to the database
// returns the DB handle
func Connect(schema string) *sql.DB {

	database_handle, error_info := sql.Open("mysql", "root:sniffer@123@/core")

	if error_info != nil {
		return nil
	}

	return database_handle
}
