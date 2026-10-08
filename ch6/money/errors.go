package money

// Error defines an error type
type Error string

// Error implements the error interface
func (e Error) Error() string {
	return string(e)
}
