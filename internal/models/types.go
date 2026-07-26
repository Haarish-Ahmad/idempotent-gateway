package models

import (
	"time"
)


type ExecState string
const(
	StateNotFound   ExecState = "NOT_FOUND"
	StateInProgress ExecState = "IN_PROGRESS"
	StateCompleted  ExecState = "COMPLETED"
)


type ContextKey string
const(
	IdompotencyKeyContextKey ContextKey = "idempotent_key"
	PayloadHashContext       ContextKey = "payload_hash"
)


type CachedResponse struct {
	StatusCode int                    `json:"status_code"`
	Headers    map[string][]string    `json:"headers"`
	Body       []byte                 `json:"body"`
	CreatedAt  time.Time              `json:"created_at"`
}


type LockResult struct {
	State          ExecState       `json:"state"`
	CachedResponse *CachedResponse `json:"cached_response,omitempty"`
}