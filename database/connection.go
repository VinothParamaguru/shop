package database

import (
	"database/sql"
	"workspace/shop/errors"

	_ "github.com/go-sql-driver/mysql"
)

// Open connect to the database
// returns the DB handle
func (db *DataBase) Open() (bool, int) {

	// open the connection to db
	dataSource := db.Config.User + ":" + db.Config.Password + "@/" + db.Config.Schema
	databaseConnection, errorInfo := sql.Open("mysql", dataSource)
	if errorInfo != nil {
		return false, errors.DbOpenFailed
	}
	db.Connector = databaseConnection
	db.ParamsMap = make(map[string]interface{})
	return true, errors.Success
}
