package models

type Modifier struct {
	ID          uint `gorm:"primaryKey"`
	OrderItemID uint
	Name        string
	UnitPrice   float64
	Quantity    int
	Amount      float64
}
