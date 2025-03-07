package model

import "time"

type Content struct {
	ID          int64      `gorm:"id"`
	UserID      int64      `gorm:"user_id"`
	CategoryID  int64      `gorm:"category_id"`
	Title       string     `gorm:"title"`
	Excerpt     string     `gorm:"excerpt"`
	Description string     `gorm:"description"`
	Image       string     `gorm:"image"`
	Tags        string     `gorm:"tags"`
	Status      string     `gorm:"status"`
	CreatedAt   time.Time  `gorm:"created_at"`
	UpdatedAt   *time.Time `gorm:"updated_at"`
	User        User       `gorm:"foreignKey:UserID"`
	Category    Category   `gorm:"foreignKey:CategoryID"`
}
