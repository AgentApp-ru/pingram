package server

import "encoding/json"

type (
	health struct {
		Space  json.RawMessage `json:"space"`
		Uptime json.RawMessage `json:"uptime"`
	}
)
