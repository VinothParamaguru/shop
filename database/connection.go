package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// connect to the database
// returns the DB handle
func (db *DataBase) Open() {

	// open the connection to db
	data_source := db.Config.User + ":" + db.Config.Password + "@/" + db.Config.Schema
	database_connection, error_info := sql.Open("mysql", data_source)
	if error_info != nil {
		panic(error_info)
	}
	db.Connector = database_connection
}
