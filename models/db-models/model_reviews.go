package dbmodels

type Review struct {
	ID string `gorm:"type:uuid;primaryKey" json:"id"`

	//foreign key - from users
	UserID string `gorm:"type:uuid" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID;references:ID"`

	//foreign key - from destinations
	DestinationID string      `gorm:"type:uuid" json:"destination_id"`
	Destination   Destination `gorm:"foreignKey:DestinationID;references:ID"`

	Rating            int      `json:"rating"`
	AdditionalComment string   `json:"additional_comment"`
	ReviewPhotos      []string `gorm:"type:text[]" json:"review_photos"`

	//foreign key - comments
	Comments []Comment `gorm:"foreignKey:ReviewID"`
}
