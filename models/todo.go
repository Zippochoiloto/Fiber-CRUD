package models

import (
	"time"
)

type Todo struct {
	ID        *string   `json:"id,omitempty" bson:"_id,omitempty"`
	Title     *string   `json:"title"`
	Completed *bool     `json:"completed"`
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time `json:"updated"`
}
