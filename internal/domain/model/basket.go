package model

import "time"

type Basket struct {
	ID        uint64    `json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Data      string    `json:"data,omitempty"`
	State     string    `json:"state,omitempty"`
	UserID    uint64    `json:"user_id"`
}
