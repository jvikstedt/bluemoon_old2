package gate

import "encoding/json"

type WorkerOut struct {
	Name    string `json:"name"`
	Payload `json:"payload"`
}

type Payload struct {
	Name    string `json:"name"`
	UserID  int    `json:"user_id"`
	Payload interface{}
}

type WorkerIn struct {
	Name    string           `json:"name"`
	Payload *json.RawMessage `json:"payload"`
}

type UserIn struct {
	Name    string           `json:"name"`
	Payload *json.RawMessage `json:"payload"`
}

type UserOut struct {
	Name    string      `json:"name"`
	Payload interface{} `json:"payload"`
}
