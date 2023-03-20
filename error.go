package fortytwo

type ErrorCode string

type Error struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
	Status  int       `json:"status"`
}

func (e *Error) Error() string {
	return e.Message
}

type RateLimitedError struct {
	Message string
}

func (e *RateLimitedError) Error() string {
	return e.Message
}
