package http

type PublisherResponse struct {
	Partition int32 `json:"partition"`
	Offset    int64 `json:"offset"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
