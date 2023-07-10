package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	apperrors "workspace/shop/errors"
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

func (db *DataBase) Select() (*sql.Rows, error) {

	// check the Select() call is valid for the SQL query with which
	// the database object was instantiated
	splittedQuery := strings.Split(db.sqlQuery, " ")
	if splittedQuery[0] != "SELECT" {
		return nil, errors.New(apperrors.DataBaseErrorDescriptions[apperrors.DbInvalidQuery])
	}
	queryForExecution := db.sqlQuery
	var arguments []interface{}
	for _, paramName := range db.Params {
		paramValue := db.ParamsMap[paramName]
		queryForExecution = strings.Replace(queryForExecution, paramName, "?", 1)
		arguments = append(arguments, paramValue)
	}
	// always use prepare statements for safety.
	statement, errorInfo := db.Connector.Prepare(queryForExecution)
	results, errorInfo := statement.Query(arguments...)
	// clear the params map
	defer db.clearParams()
	if errorInfo != nil {
		return nil, errors.New(apperrors.DataBaseErrorDescriptions[apperrors.DbErrorSelectQueryExecution])
	}
	return results, nil
}

func (db *DataBase) Execute() error {

	queryForExecution := db.sqlQuery
	var arguments []interface{}
	for _, paramName := range db.Params {
		paramValue := db.ParamsMap[paramName]
		queryForExecution = strings.Replace(queryForExecution, paramName, "?", 1)
		arguments = append(arguments, paramValue)
	}

	// always use prepare statements for safety.
	statement, errorInfo := db.Connector.Prepare(queryForExecution)
	_, errorInfo = statement.Exec(arguments...)
	// clear the params map
	defer db.clearParams()
	if errorInfo != nil {
		fmt.Println(errorInfo.Error())
		return errors.New(apperrors.DataBaseErrorDescriptions[apperrors.DbErrorQueryExecution])
	}
	return nil
}

func (db *DataBase) InitSql(sql string) {
	// store the query for later use
	db.sqlQuery = sql
	db.ParamsMap = make(map[string]interface{})
	sqlParams := utilities.GetSqlParams(sql)
	for _, param := range sqlParams {
		db.ParamsMap[param] = 0
		db.Params = append(db.Params, param)
	}
}

func (db *DataBase) BindParam(param string, value any) error {

	// Bind has been called without any params for the prepared statement
	if len(db.ParamsMap) == 0 {
		return errors.New(apperrors.DataBaseErrorDescriptions[apperrors.DbBindParamNotApplicable])
	}

	db.ParamsMap[param] = value

	return nil
}

func (db *DataBase) clearParams() {

	// clear the entries in the map that holds the params
	// this is to ensure the params are re-initiated when the Select
	// call is used multiple times from the same connection
	db.Params = nil
	for key, _ := range db.ParamsMap {
		delete(db.ParamsMap, key)
	}
}
