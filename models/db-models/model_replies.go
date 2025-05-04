package dbmodels

type Reply struct {
	ID string `gorm:"type:uuid;primaryKey" json:"id"`

	//foreign key - from comments
	CommentID string  `gorm:"type:uuid" json:"comment_id"`
	Comment   Comment `gorm:"foreignKey:CommentID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`

	//foreign key - from users
	UserID string `gorm:"type:uuid" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`

	Reply string `json:"reply"`
}