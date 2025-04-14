package models

import "time"

type Order struct {
	ID             uint      `gorm:"primaryKey"`
	RestaurantID   uint
	TableNumber    string
	SquareOrderID  string
	OpenedAt       time.Time
	IsClosed       bool
	OrderItems     []OrderItem
	Payments       []Payment
}
