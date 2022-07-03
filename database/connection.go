package database

import (
	"database/sql"
	app_error "workspace/shop/error"

	_ "github.com/go-sql-driver/mysql"
)

// Open connect to the database
// returns the DB handle
func (db *DataBase) Open() (bool, int) {

	// open the connection to db
	dataSource := db.Config.User + ":" + db.Config.Password + "@/" + db.Config.Schema
	databaseConnection, errorInfo := sql.Open("mysql", dataSource)
	if errorInfo != nil {
		return false, app_error.DbOpenFailed
	}
	db.Connector = databaseConnection
	return true, app_error.Success
}
