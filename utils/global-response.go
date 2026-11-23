package utils

import "github.com/michael870917/gocourse_meta/meta"

type GlobalResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
	Meta   *meta.Meta  `json:"meta,omitempty"`
}
