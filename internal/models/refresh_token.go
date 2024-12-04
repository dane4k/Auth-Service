package models

import "time"

type RefreshToken struct {
	Id          string    `gorm:"type:uuid;primaryKey"`
	UserId      string    `gorm:"type:uuid;not null"`
	TokenHashed string    `gorm:"type:text;not null"`
	UserIp      string    `gorm:"type:text;not null"`
	Expires     time.Time `gorm:"type:timestamp;not null"`
	CreatedAt   time.Time `gorm:"type:timestamp;autoCreateTime"`
}
