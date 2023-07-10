package errors

// Error codes - http specific
const (
	HttpReadFailed         = 4000
	HttpUnknownContentType = 4001
)

// HttpErrorDescriptions and mappings
var HttpErrorDescriptions = map[int]string{
	HttpReadFailed:         "Unable to read the http request",
	HttpUnknownContentType: "Unknown content type in the request",
}
