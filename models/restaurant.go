package models

type Restaurant struct {
	ID                uint   `gorm:"primaryKey"`
	Name              string
	SquareAccessToken string
	LocationID        string
}
