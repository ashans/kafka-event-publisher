package http

import "ev_pub/internal/common"

type PublisherRequest struct {
	Topic       string                    `json:"topic"`
	Key         common.TypeValWithOptions `json:"key"`
	Value       common.TypeValWithOptions `json:"value"`
	Partitioner common.TypeValWithOptions `json:"partitioner"`
	Headers     map[string]string         `json:"headers"`
}
