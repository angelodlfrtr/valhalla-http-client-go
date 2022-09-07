package client

// ErrorResponse from the valhalla server
type ErrorResponse struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error"`
	StatusCode   int    `json:"status_code"`
	Status       string `json:"status"`
}

// Error as string
func (err *ErrorResponse) Error() string {
	return err.Status + ": " + err.ErrorMessage
}
