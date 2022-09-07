package client

// Point define a geographical point
type Point struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

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
