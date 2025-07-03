package main

import "time"

type OrderItem struct {
	ID          int32   `gorm:"column:id"`
	OrderID     *int32  `gorm:"column:order_id"`
	ProductName string  `gorm:"column:product_name"`
	Price       float64 `gorm:"column:price"`
	Quantity    *int32  `gorm:"column:quantity"`
}

type OrderWithItems struct {
	ID           int32       `gorm:"column:id"`
	CustomerName string      `gorm:"column:customer_name"`
	CreatedAt    *time.Time  `gorm:"column:created_at"`
	Itens        []OrderItem `gorm:"foreignKey:OrderID;references:ID"`
}

func (OrderWithItems) TableName() string {
	return "orders"
}
