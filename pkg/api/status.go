package api

// Code represents an API status code.
type Code int

const (
	// OK is returned on success.
	//
	// HTTP Mapping: 200 OK.
	OK Code = iota

	// Unknown error.
	//
	// HTTP Mapping: 500 Internal Server Error.
	Unknown

	// InvalidArgument indicates client specified an invalid argument.
	//
	// HTTP Mapping: 400 Bad Request.
	InvalidArgument

	// ...
)
