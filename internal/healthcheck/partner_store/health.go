package partner_store

import "encoding/json"

type (
	health struct {
		Health map[string]json.RawMessage `json:"health"`
		Stats  map[string]json.RawMessage `json:"stats"`
	}
)
