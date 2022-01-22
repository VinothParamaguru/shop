package database

import (
	"database/sql"
	app_error "workspace/shop/error"

	_ "github.com/go-sql-driver/mysql"
)

// connect to the database
// returns the DB handle
func (db *DataBase) Open() (bool, int) {

	// open the connection to db
	data_source := db.Config.User + ":" + db.Config.Password + "@/" + db.Config.Schema
	database_connection, error_info := sql.Open("mysql", data_source)
	if error_info != nil {
		return false, app_error.DbOpenFailed
	}
	db.Connector = database_connection
	return true, app_error.Success
}
