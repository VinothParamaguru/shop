package errors

// Error codes - database specific
const (
	DbOpenFailed             = 2000
	DbBindParamNotApplicable = 2001
	DbInvalidQuery           = 2003
	DbErrorQueryExecution    = 2004
)

// DataBaseErrorDescriptions and mappings
var DataBaseErrorDescriptions = map[int]string{
	DbOpenFailed:             "Error opening the database",
	DbBindParamNotApplicable: "Database error, bind param not applicable",
	DbInvalidQuery:           "Database error, invalid sql query",
	DbErrorQueryExecution:    "Database error, error executing query",
}
