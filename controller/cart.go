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
		routes.POST("/", func(c *gin.Context) {
			GetUserByEmail(c, db)
		})
		routes.PUT("/update", func(c *gin.Context) {
			UpdateUser(c, db)
		})
		routes.PUT("/changepwd", func(c *gin.Context) {
			UpdatePassword(c, db)
		})
	}
}

func GetCart(c *gin.Context, db *gorm.DB) {
	// ใช้ DTO สำหรับรับ JSON
	var request dto.UserDataRequest

	// ตรวจสอบ JSON ว่าถูกต้องหรือไม่
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// ตรวจสอบว่า email ไม่เป็น nil
	if request.Email == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	// ค้นหาผู้ใช้ตามอีเมล
	var user model.Customer
	result := db.Where("email = ?", request.Email).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	// ส่งข้อมูลผู้ใช้กลับไป
	c.JSON(http.StatusOK, gin.H{
		"message": "User found",
		"user":    user,
	})
}
