package errors

// Error codes - database specific
const (
	DbOpenFailed                = 2000
	DbBindParamNotApplicable    = 2001
	DbInvalidQuery              = 2003
	DbErrorSelectQueryExecution = 2004
	DbErrorQueryExecution       = 2005
	DbErrorScanning             = 2006
	DbErrorCreatingPreparedStmt = 2007
)

// DataBaseErrorDescriptions and mappings
var DataBaseErrorDescriptions = map[int]string{
	DbOpenFailed:                "Error opening the database",
	DbBindParamNotApplicable:    "Database error, bind param not applicable",
	DbInvalidQuery:              "Database error, invalid sql query",
	DbErrorSelectQueryExecution: "Database error, error executing select query",
	DbErrorQueryExecution:       "Database error, error executing query",
	DbErrorScanning:             "Database error, error executing query",
	DbErrorCreatingPreparedStmt: "Database error, error in creating prepared statement",
}
