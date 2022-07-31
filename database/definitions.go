package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DataBaseConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Schema   string `json:"schema"`
}

type DataBase struct {
	Connector *sql.DB
	Config    DataBaseConfig
	Params    []string
	ParamsMap map[string]interface{}
	sqlQuery  string
}

type Field struct {
	Name  string
	Value interface{}
}
