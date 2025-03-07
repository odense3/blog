package model

import "time"

type Category struct {
	ID        int64      `gorm:"id"`
	UserID    int64      `gorm:"user_id"`
	Title     string     `gorm:"title"`
	Slug      string     `gorm:"slug"`
	CreatedAt time.Time  `gorm:"created_at"`
	UpdatedAt *time.Time `gorm:"updated_at"`
	User      User       `gorm:"foreignKey:UserID"`
}
