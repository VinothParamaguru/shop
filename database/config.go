package database

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetDataBaseConfig() DataBaseConfig {

	// not the best solution to identify the execution environment is cloud-run
	// but assuming we pass the credentials from secret manager to the environment
	// variables, this should work
	// please ensure the environment vars are populated correct before executing the code

	var databaseConfig DataBaseConfig

	dbUserName := os.Getenv("db-username")
	dbPassword := os.Getenv("db-password")
	dbSchema := os.Getenv("db-schema")

	if dbUserName != "" && dbPassword != "" && dbSchema != "" {
		// cloud-run
		databaseConfig.User = dbUserName
		databaseConfig.Password = dbPassword
		databaseConfig.Schema = dbSchema
	} else {
		// local/vm/possible k8s where mounting is viable
		// read the database configuration file
		// using relative path for now. This should be replaced by absolute path of the file
		configFileRelativePath := "config/database_config.json"
		configFileAbsolutePath, errorInfo := filepath.Abs(configFileRelativePath)
		if errorInfo != nil {
			panic(errorInfo)
		}
		// read
		data, errorInfo := ioutil.ReadFile(configFileAbsolutePath)
		if errorInfo != nil {
			panic(errorInfo)
		}
		// unmarshal
		errorInfo = json.Unmarshal(data, &databaseConfig)
		if errorInfo != nil {
			panic(errorInfo)
		}
	}
	return databaseConfig
}
