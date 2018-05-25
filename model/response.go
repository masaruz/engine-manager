package model

import (
	"encoding/json"
	"net/http"
)

// Response when send back to client
type Response struct {
	StatusCode StatusCode `json:"statusCode"`
	Message    string     `json:"message,omitempty"`
	Logs       []string   `json:"logs,omitempty"`
}

// StatusCode for response's status
type StatusCode int

// StatusCode number
const (
	StatusOK             StatusCode = 200
	StatusGetPodFail     StatusCode = 300
	StatusGetPodLogsFail StatusCode = 301
)

// Reply payload to sender
func Reply(w http.ResponseWriter, payload Response) {
	j, _ := json.Marshal(payload)
	w.Write(j)
}
