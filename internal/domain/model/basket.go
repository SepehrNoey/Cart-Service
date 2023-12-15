package model

import "time"

type Basket struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      []byte
	State     string
}
