package utils

import (
	"fmt"
	"time"

	dbmodels "github.com/mhmdfathan/funsies-be/models/db-models"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func StartCron(db *gorm.DB) {
	fmt.Println("Cron job started")
	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("0 * * * * *", func() {
		CheckPendingUsers(db)
	})

	if err != nil {
		fmt.Printf("Failed to schedule cron job: %v", err)
	}

	c.Start()
}

func CheckPendingUsers(db *gorm.DB) {
	fmt.Println("Checking pending users")

	var users []dbmodels.User
	if err := db.Where("is_active = ? AND created_at < ?", false, time.Now().Add(-24*time.Hour)).Find(&users).Error; err != nil {
		fmt.Println("Error querying pending users:", err)
		return
	}

	for _, user := range users {
		fmt.Printf("Deleting pending user: %s (created: %s)", user.Email, user.CreatedAt)
		
		if err := db.Delete(&dbmodels.User{}, "id = ?", user.ID).Error; err != nil {
			fmt.Printf("Failed to delete user %s: %v", user.ID, err)
		}
	}
}
