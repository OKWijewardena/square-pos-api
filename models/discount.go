package models

type Discount struct {
	ID          uint `gorm:"primaryKey"`
	OrderItemID uint
	Name        string
	IsPercentage bool
	Value       float64
	Amount      float64
}
