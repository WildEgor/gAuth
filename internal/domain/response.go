package domains

// ErrorResponseDomain
type ErrorResponseDomain struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

// VoidResponseDomain
type VoidResponseDomain struct {
	Message string `json:"message"`
}

func NewVoidResponseDomain() *VoidResponseDomain {
	return &VoidResponseDomain{
		Message: "success",
	}
}
