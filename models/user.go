package models

type User struct {
	ID       *string `json:"id,omitempty" bson:"_id,omitempty"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type ChangePassword struct {
	Email       *string `json:"email"`
	OldPassword *string `json:"oldPassword"`
	NewPassword *string `json:"newPassword"`
}
