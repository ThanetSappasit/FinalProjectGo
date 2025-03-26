package controller

import (
	"final_go/dto"
	"final_go/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CountryController sets up country-related routes
func ProductController(router *gin.Engine, db *gorm.DB) {
	routes := router.Group("/product")
	{
		routes.GET("/", func(c *gin.Context) {
			GetAllProduct(c, db)
		})
		routes.POST("/search", func(c *gin.Context) {
			SearchProduct(c, db)
		})
	}
}

func GetAllProduct(c *gin.Context, db *gorm.DB) {
	var product []model.Product
	result := db.Find(&product)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func SearchProduct(c *gin.Context, db *gorm.DB) {
	var req dto.ProductSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var products []model.Product
	query := db.Model(&model.Product{})

	// Search by product name (case-insensitive, partial match)
	if req.ProductName != nil && *req.ProductName != "" {
		query = query.Where("LOWER(product_name) LIKE LOWER(?)", "%"+*req.ProductName+"%")
	}

	// Search by description (case-insensitive, partial match)
	if req.Description != nil && *req.Description != "" {
		query = query.Where("LOWER(description) LIKE LOWER(?)", "%"+*req.Description+"%")
	}

	// Price range filtering
	if req.MinPrice != nil && *req.MinPrice != "" {
		query = query.Where("CAST(price AS DECIMAL) >= ?", *req.MinPrice)
	}

	if req.MaxPrice != nil && *req.MaxPrice != "" {
		query = query.Where("CAST(price AS DECIMAL) <= ?", *req.MaxPrice)
	}

	// Execute the query
	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if any products were found
	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No products found matching the search criteria",
			"count":   0,
		})
		return
	}

	// Return search results with additional metadata
	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"count":    len(products),
	})
}
