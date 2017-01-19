package worker

import "encoding/json"

type GateIn struct {
	Name    string           `json:"name"`
	Payload *json.RawMessage `json:"payload"`
}
