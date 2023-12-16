package model

import "time"

type Basket struct {
	ID        uint64    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Data      string    `json:"data,omitempty"`
	State     string    `json:"state,omitempty"`
}
