package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"workspace/shop/errors"
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

func (db *DataBase) Select() (bool, int, *sql.Rows) {

	// check the Select() call is valid for the SQL query with which
	// the database object was instantiated
	splittedQuery := strings.Split(db.sqlQuery, " ")
	if splittedQuery[0] != "SELECT" {
		return false, errors.DbInvalidQuery, nil
	}
	queryForExecution := db.sqlQuery
	var arguments []interface{}
	for key, value := range db.ParamsMap {
		queryForExecution = strings.Replace(queryForExecution, key, "?", 1)
		arguments = append(arguments, value)
	}

	// always use prepare statements for safety.
	statement, errorInfo := db.Connector.Prepare(queryForExecution)

	results, errorInfo := statement.Query(arguments...)

	// clear the params map
	defer db.clearParams()

	if errorInfo != nil {
		return false, errors.DbErrorQueryExecution, nil
	}

	return true, errors.Success, results
}

func (db *DataBase) Execute() (bool, int) {

	queryForExecution := db.sqlQuery
	var arguments []interface{}
	for key, value := range db.ParamsMap {
		queryForExecution = strings.Replace(queryForExecution, key, "?", 1)
		arguments = append(arguments, value)
	}

	// always use prepare statements for safety.
	statement, errorInfo := db.Connector.Prepare(queryForExecution)

	_, errorInfo = statement.Exec(arguments...)

	// clear the params map
	defer db.clearParams()

	if errorInfo != nil {
		return false, errors.DbErrorQueryExecution
	}

	return true, errors.Success
}

func (db *DataBase) InitSql(sql string) {
	// store the query for later use
	db.sqlQuery = sql
	db.ParamsMap = make(map[string]interface{})
	sqlParams := utilities.GetSqlParams(sql)
	for _, param := range sqlParams {
		db.ParamsMap[param] = 0
	}
}

func (db *DataBase) BindParam(param string, value any) (bool, int) {

	// Bind has been called without any params for the prepared statement
	if len(db.ParamsMap) == 0 {
		return false, errors.DbBindParamNotApplicable
	}

	db.ParamsMap[param] = value

	return true, errors.Success
}

func (db *DataBase) clearParams() {

	// clear the entries in the map that holds the params
	// this is to ensure the params are re-initiated when the Select
	// call is used multiple times from the same connection
	for key, _ := range db.ParamsMap {
		delete(db.ParamsMap, key)
	}
}
