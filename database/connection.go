package database

import (
	"database/sql"
	"errors"
	apperrors "workspace/shop/errors"

	_ "github.com/go-sql-driver/mysql"
)

// Open connect to the database
// returns the DB handle
func (db *DataBase) Open() error {

	// open the connection to db
	dataSource := db.Config.User + ":" + db.Config.Password + "@/" + db.Config.Schema
	databaseConnection, errorInfo := sql.Open("mysql", dataSource)
	if errorInfo != nil {
		return errors.New(apperrors.DataBaseErrorDescriptions[apperrors.DbOpenFailed])
	}
	db.Connector = databaseConnection
	return nil
}
