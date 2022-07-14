package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	app_error "workspace/shop/error"
	"workspace/shop/utilities"
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

func (db *DataBase) InitSql(sql string) {

	sqlParams := utilities.GetSqlParams(sql)
	fmt.Print(sqlParams)
	for _, param := range sqlParams {
		db.ParamsMap[param] = 0
		fmt.Print(param)
	}
}

func (db *DataBase) BindParam(param string, value any) (bool, int) {

	// Bind has been called without any params for the prepared statement
	if len(db.ParamsMap) == 0 {
		return false, app_error.DbBindParamNotApplicable
	}

	db.ParamsMap[param] = value

	return true, app_error.Success
}
