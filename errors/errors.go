package errors

import "errors"

// Success code common for all the modules
const (
	Success = 0
)

var ErrorTableUnified = []map[int]string{
	SecurityErrorDescriptions,
	ApplicationErrorDescriptions,
	DataBaseErrorDescriptions,
	HttpErrorDescriptions,
}

func GetError(errorCode int) error {
	for _, errorTable := range ErrorTableUnified {
		errorDescription, found := errorTable[errorCode]
		if found {
			return errors.New(errorDescription)
		}
	}
	return errors.New("unknown error")
}
