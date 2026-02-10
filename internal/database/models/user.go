package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Login    string `gorm:"unique;not null" json:"login"`
	Password string `gorm:"not null" json:"-"`
}
