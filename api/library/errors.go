package library

// ApiError holds errors to show to the user.
type ApiError struct {
	Code int
	Msg  string
}

var _ error = (*ApiError)(nil)

func NewApiError(code int, msg string) *ApiError {
	return &ApiError{
		Code: code,
		Msg:  msg,
	}
}

func (e *ApiError) Error() string {
	return e.Msg
}

var (
	Err404 = NewApiError(404, "Not Found")
	Err500 = NewApiError(500, "Internal Server Error")
)
