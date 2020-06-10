package utils

import (
	"encoding/json"
	"github.com/greenac/artemis/logger"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Payload interface{} `json:"payload"`
	Message string      `json:"message"`
}

func (r *Response) Respond(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	if r.Code != http.StatusOK {
		w.WriteHeader(r.Code)
	}

	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		logger.Error("Response::Respond Failed to encode json", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (r *Response) SetPayload(key string, payload interface{}) {
	r.Payload = map[string]interface{}{key: payload}
}

func (r *Response) SetPayloadNoKey(payload interface{}) {
	r.Payload = payload
}
