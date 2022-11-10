package application

import "fmt"

// ErrorIDFormat is used when the data repository cannot use the bookInfo's ID is not in an acceptable format the the underlying provider
type ErrorIDFormat struct {
	ID string
}

func (e *ErrorIDFormat) Error() string {
	return fmt.Sprintf("ID is not in expected format %s", e.ID)
}

// ErrorCannotFindBookStock is used when the bookInfo with the given ID cannot be found on the underlying data source
type ErrorCannotFindBookStock struct {
	ID string
}

func (e *ErrorCannotFindBookStock) Error() string {
	return fmt.Sprintf("Cannot find the bookInfo with the ID %s", e.ID)
}

// ErrorParsePayload is used when the payload is cannot be parsed by the communications package
type ErrorParsePayload struct{}

func (e *ErrorParsePayload) Error() string {
	return "Cannot parse payload"
}

// ErrorReadPayload is used when the payload is cannot be read by the communications package
type ErrorReadPayload struct{}

func (e *ErrorReadPayload) Error() string {
	return "Cannot read payload"
}

// ErrorPayloadMissing is used when the communication package expects a payload and there is none
type ErrorPayloadMissing struct{}

func (e *ErrorPayloadMissing) Error() string {
	return "Payload is missing"
}
