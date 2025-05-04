package dbmodels

import (
	"time"
)

type User struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string    `gorm:"type:varchar(100);unique" json:"email"`
	Username  string    `gorm:"type:varchar(50);unique" json:"username"`
	Password  string    `json:"password"`
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
 	BirthDate time.Time `json:"birth_date"`
	Gender    uint8     `json:"gender"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Session   string    `json:"session"`
	IsActive  bool 		`gorm:"default:false" json:"is_active"`
	Otp		  string 	`json:"otp"`

	//foreign key - reviews
	Reviews []Review 	`gorm:"foreignKey:UserID"`

	//foreign key - comments
	Comments []Comment `gorm:"foreignKey:UserID"`

	//foreign key - activation tokens
	ActivationTokens []ActivationToken `gorm:"foreignKey:UserID"`
}
