package models

import (
	"time"
)

// Feedback represents user feedback submission in the system
type Feedback struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	DeviceID  string    `gorm:"not null;size:255;index" json:"device_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for Feedback model
func (Feedback) TableName() string {
	return "feedbacks"
}
