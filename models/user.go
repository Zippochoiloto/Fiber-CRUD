package models

type User struct {
	ID       *string `json:"id,omitempty" bson:"_id,omitempty"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}
