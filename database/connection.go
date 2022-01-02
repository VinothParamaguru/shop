package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func InitStatement(sql_query string) {

	database_handle, error_info := sql.Open("mysql", "root:sniffer@123@/core")
	if error_info != nil {
		panic(error_info.Error())
	}
	return database_handle
}
