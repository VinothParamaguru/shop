package database

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

func GetDataBaseConfig() DataBaseConfig {

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
	var databaseConfig DataBaseConfig
	errorInfo = json.Unmarshal(data, &databaseConfig)
	if errorInfo != nil {
		panic(errorInfo)
	}

	return databaseConfig
}
