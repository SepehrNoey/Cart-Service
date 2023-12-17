package model

type User struct {
	ID       uint64 `json:"user_id,omitempty" gorm:"primaryKey"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
