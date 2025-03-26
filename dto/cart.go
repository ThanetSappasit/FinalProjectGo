package dto

type AddToCartRequest struct {
	CustomerID int    `json:"customer_id" binding:"required"`
	CartName   string `json:"cart_name" binding:"required"`
	ProductID  int    `json:"product_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
}

type CartResponse struct {
	CartID     int                `json:"cart_id"`
	CustomerID int                `json:"customer_id"`
	CartName   string             `json:"cart_name"`
	CartItems  []CartItemResponse `json:"cart_items"`
	TotalItems int                `json:"total_items"`
}

type CartItemResponse struct {
	CartItemID  int    `json:"cart_item_id"`
	ProductID   int    `json:"product_id"`
	Quantity    int    `json:"quantity"`
	ProductName string `json:"product_name"`
	Price       string `json:"price"`
}
