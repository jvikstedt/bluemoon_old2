package worker

import "encoding/json"

type GateIn struct {
	Name    string           `json:"name"`
	Payload *json.RawMessage `json:"payload"`
}

type GateOut struct {
	Name    string      `json:"name"`
	Payload interface{} `json:"payload"`
}

type ToUsers struct {
	UserIds []int       `json:"user_ids"`
	Payload interface{} `json:"payload"`
}

type Action interface {
	Run(room Room, delta float64)
}
