package models

type OrderItem struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint
	Name      string
	UnitPrice float64
	Quantity  int
	Amount    float64
	Discounts []Discount
	Modifiers []Modifier
}
