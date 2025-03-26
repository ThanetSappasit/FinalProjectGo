package controller

import (
	"final_go/dto"
	"final_go/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CountryController sets up country-related routes
func CartController(router *gin.Engine, db *gorm.DB) {
	routes := router.Group("/cart")
	{
		routes.POST("/add", func(c *gin.Context) {
			AddItemToCart(c, db)
		})
		// routes.PUT("/update", func(c *gin.Context) {
		// 	UpdateUser(c, db)
		// })
		// routes.PUT("/changepwd", func(c *gin.Context) {
		// 	UpdatePassword(c, db)
		// })
	}
}

// เพิ่มสินค้าลงในรถเข็นตามชื่อของรถเข็นที่ต้องการ
func AddItemToCart(c *gin.Context, db *gorm.DB) {
	// Parse the request body
	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a database transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find or create cart
	var cart model.Cart
	result := tx.Where("customer_id = ? AND cart_name = ?", req.CustomerID, req.CartName).First(&cart)

	// If cart doesn't exist, create a new one
	if result.Error != nil {
		cart = model.Cart{
			CustomerID: req.CustomerID,
			CartName:   req.CartName,
		}
		if err := tx.Create(&cart).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
			return
		}
	}

	// Check if the product exists
	var product model.Product
	if err := tx.First(&product, req.ProductID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Check if the product is already in the cart
	var existingCartItem model.CartItem
	result = tx.Where("cart_id = ? AND product_id = ?", cart.CartID, req.ProductID).First(&existingCartItem)

	if result.Error == nil {
		// Product exists in cart, update quantity
		existingCartItem.Quantity += req.Quantity
		if err := tx.Save(&existingCartItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
			return
		}
	} else {
		// Product not in cart, create new cart item
		cartItem := model.CartItem{
			CartID:    cart.CartID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		if err := tx.Create(&cartItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	// Prepare the response
	var cartItems []dto.CartItemResponse
	var totalItems int

	// Fetch cart items with product details
	if err := db.Table("cart_item").
		Select("cart_item.cart_item_id, cart_item.product_id, cart_item.quantity, product.product_name, product.price").
		Joins("JOIN product ON cart_item.product_id = product.product_id").
		Where("cart_item.cart_id = ?", cart.CartID).
		Scan(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart items"})
		return
	}

	// Calculate total items
	for _, item := range cartItems {
		totalItems += item.Quantity
	}

	// Prepare full cart response
	response := dto.CartResponse{
		CartID:     cart.CartID,
		CustomerID: cart.CustomerID,
		CartName:   cart.CartName,
		CartItems:  cartItems,
		TotalItems: totalItems,
	}

	c.JSON(http.StatusOK, response)
}
