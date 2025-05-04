package dbmodels

type Comment struct {
	ID string `gorm:"type:uuid;primaryKey" json:"id"`

	//foreign key - from reviews
	ReviewID string `gorm:"type:uuid" json:"review_id"`
	Review   Review `gorm:"foreignKey:ReviewID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`

	//foreign key - from users
	UserID string `gorm:"type:uuid" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`

	Comment string `json:"comment"`

	//foreign key - replies
	// Replies []Reply `gorm:"foreignKey:CommentID"`
}