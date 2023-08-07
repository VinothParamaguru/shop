package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	apperrors "shop/errors"

	_ "github.com/go-sql-driver/mysql"
)

// Open connect to the database
// returns the DB handle
func (db *DataBase) Open() error {

	var dataSource string
	// unix socket path to be used under cloud-run environment
	unixSocketPath := os.Getenv("db-unix-socket-path")
	if unixSocketPath != "" {
		dataSource = fmt.Sprintf("%s:%s@unix(%s)/%s?parseTime=true",
			db.Config.User, db.Config.Password, unixSocketPath, db.Config.Schema)
	} else {
		dataSource = db.Config.User + ":" + db.Config.Password + "@/" + db.Config.Schema
	}
	// open the connection to db
	databaseConnection, errorInfo := sql.Open("mysql", dataSource)
	if errorInfo != nil {
		return errors.New(apperrors.DataBaseErrorDescriptions[apperrors.DbOpenFailed])
	}
	db.Connector = databaseConnection
	return nil
}

func (db *DataBase) Close() {
	_ = db.Connector.Close()
}
