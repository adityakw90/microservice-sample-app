package params

// PaginationParam represents common pagination parameters.
type PaginationParam struct {
	Page  int
	Limit int
}

// ListParam represents common list query parameters.
type ListParam struct {
	PaginationParam
}
