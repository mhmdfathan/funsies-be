package dbmodels

type FollowingFollowed struct {
	ID string `gorm:"type:uuid;primaryKey" json:"id"`

	// Foreign key - the user who follows
	FollowingID   string `gorm:"type:uuid" json:"following_id"`
	FollowingUser User   `gorm:"foreignKey:FollowingID;references:ID" json:"following_user"`

	// Foreign key - the user being followed
	FollowedID   string `gorm:"type:uuid" json:"followed_id"`
	FollowedUser User   `gorm:"foreignKey:FollowedID;references:ID" json:"followed_user"`
}
