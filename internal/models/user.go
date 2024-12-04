package models

type User struct {
	Id    string `gorm:"type:uuid;primaryKey" json:"id"`
	Email string `gorm:"type:text;unique;not null" json:"email"`
}
