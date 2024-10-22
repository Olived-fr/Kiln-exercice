package api

// Error represents an API error.
type Error struct {
	Code    Code
	Err     error
	Message string
}

func (e Error) Error() string {
	return e.Err.Error()
}

func NewError(code Code, message string, err error) *Error {
	return &Error{
		Code:    code,
		Err:     err,
		Message: message,
	}
}
