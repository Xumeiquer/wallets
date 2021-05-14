package models

import "encoding/json"

type APIResponse struct {
	Type string          `json:"type"`
	Msg  json.RawMessage `json:"msg"`
}
