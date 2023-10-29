package domains

// ErrorResponseDomain
type ErrorResponseDomain struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}
