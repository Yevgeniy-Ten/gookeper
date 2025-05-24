package models

import "time"

type Secret struct {
	ID        int64       `json:"id"`
	UserID    int64       `json:"user_id"`
	Type      string      `json:"type"`
	Data      []byte      `json:"data"`
	Meta      interface{} `json:"meta"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type CreateSecretDTO struct {
	Type string      `json:"type" binding:"required,oneof=password text binary card"`
	Data []byte      `json:"data" binding:"required"`
	Meta interface{} `json:"meta"`
}

type UpdateSecretDTO struct {
	Type string      `json:"type" binding:"omitempty,oneof=password text binary card"`
	Data []byte      `json:"data" binding:"omitempty"`
	Meta interface{} `json:"meta"`
}
