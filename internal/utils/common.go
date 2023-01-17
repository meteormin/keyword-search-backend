package utils

type StatusResponse struct {
	Status bool `json:"status"`
}

type DataResponse[T interface{}] struct {
	Data T `json:"data"`
}
