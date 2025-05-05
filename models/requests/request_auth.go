package requests

import "time"

type RequestRegister struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required,min=8"`
	Phone     string `json:"phone" validate:"required,phoneval"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	BirthDate time.Time `json:"birth_date" validate:"required"`
	Gender 		uint8 `json:"gender" validate:"required,number"`
}

type PendingUser struct {
	ID 			string
	IsActive 	bool
	Token 		string
	ExpiresAt	time.Time
}