package models

import (
	"encoding/json"
	"net/http"
	"time"
)

type Error struct {
	StatusCode int
	Reason     string
	DateTime   time.Time
}

type CacheData struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Timestamp time.Time
	Prev      *CacheData
	Next      *CacheData
}

type Success struct {
	Data interface{} `json:"data"`
}

// Error returns a string representation of the error.
// It returns the reason for the error stored in the Error struct.
func (e Error) Error() string {
	return e.Reason
}

// SetError sets the HTTP status code of the response, writes the error details to the response body in JSON format, and sends the response.
// It takes the response writer, the HTTP status code, and a reason for the error as parameters.
func SetError(w http.ResponseWriter, statusCode int, reason string) {
	w.WriteHeader(statusCode)

	e := Error{StatusCode: statusCode, Reason: reason, DateTime: time.Now()}

	json.NewEncoder(w).Encode(e)
}
