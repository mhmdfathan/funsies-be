package dbmodels

import "time"

type ActivationToken struct {
	ID        string    `gorm:"type:uuid;primaryKey"`

	//foreign key - from users
	UserID    string    `gorm:"type:uuid;not null"`
	User      User 		`gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	Token     string    `gorm:"type:varchar(255);not null;unique"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}
