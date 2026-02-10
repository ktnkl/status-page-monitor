package models

type Server struct {
	ID          uint   `gorm:"primaryKey"`
	Url         string `gorm:"type:varchar(255);not null"`
	Interval    int    `gorm:"not null"`
	Checkedat   string `gorm:"type:varchar(29);not null"`
	Nextcheckat string `gorm:"type:varchar(29);not null"`
	Status      string `gorm:"default:'inactive';not null"`
}
