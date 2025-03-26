// Code generated by sql2gorm. DO NOT EDIT.
package model

import (
	"time"
)

type Cart struct {
	CartID     int       `gorm:"column:cart_id;primary_key;AUTO_INCREMENT"`
	CustomerID int       `gorm:"column:customer_id;NOT NULL"`
	CartName   string    `gorm:"column:cart_name"`
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

func (m *Cart) TableName() string {
	return "cart"
}

type CartItem struct {
	CartItemID int       `gorm:"column:cart_item_id;primary_key;AUTO_INCREMENT"`
	CartID     int       `gorm:"column:cart_id;NOT NULL"`
	ProductID  int       `gorm:"column:product_id;NOT NULL"`
	Quantity   int       `gorm:"column:quantity;NOT NULL"`
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	Cart 		Cart `gorm:"foreignKey:CartID;references:CartID"`
}

func (m *CartItem) TableName() string {
	return "cart_item"
}