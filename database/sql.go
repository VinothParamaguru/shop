package database

import (
	_ "github.com/go-sql-driver/mysql"
)

func (db *DataBase) Insert(table_name string, fields []Field) {

	//total_fields := len(fields)
	insert_query := "INSERT INTO " + table_name + " ("
	for _, field := range fields {
		insert_query += field.Name
		insert_query += ","
	}
	insert_query = insert_query[:len(insert_query)-1]
	insert_query += ") VALUES ("
	for i := 0; i < len(fields); i++ {
		insert_query += "?,"
	}
	insert_query = insert_query[:len(insert_query)-1]
	insert_query += ")"

	// prepare statement for insert

	insert_statement, error_info := db.Connector.Prepare(insert_query)

	if error_info != nil {
		panic(error_info)
	}

	var values []interface{}

	for _, field := range fields {
		values = append(values, field.Value)
	}

	_, error_info = insert_statement.Exec(values...)

	if error_info != nil {
		panic(error_info)
	}
}
