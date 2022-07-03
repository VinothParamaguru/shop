package database

import (
	_ "github.com/go-sql-driver/mysql"
)

func (db *DataBase) Insert(tableName string, fields []Field) {

	//total_fields := len(fields)
	insertQuery := "INSERT INTO " + tableName + " ("
	for _, field := range fields {
		insertQuery += field.Name
		insertQuery += ","
	}
	insertQuery = insertQuery[:len(insertQuery)-1]
	insertQuery += ") VALUES ("
	for i := 0; i < len(fields); i++ {
		insertQuery += "?,"
	}
	insertQuery = insertQuery[:len(insertQuery)-1]
	insertQuery += ")"

	// prepare statement for insert

	insertStatement, errorInfo := db.Connector.Prepare(insertQuery)

	if errorInfo != nil {
		panic(errorInfo)
	}

	var values []interface{}

	for _, field := range fields {
		values = append(values, field.Value)
	}

	_, errorInfo = insertStatement.Exec(values...)

	if errorInfo != nil {
		panic(errorInfo)
	}
}
