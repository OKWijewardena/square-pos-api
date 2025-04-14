package models

type Payment struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint
	BillAmount float64
	TipAmount  float64
	PaymentID  string
}
