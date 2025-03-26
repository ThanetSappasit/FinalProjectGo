package controller

import (
	"final_go/dto"
	"final_go/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CountryController sets up country-related routes
func LoginController(router *gin.Engine, db *gorm.DB) {
	routes := router.Group("/auth")
	{
		// routes.GET("/", func(c *gin.Context) {
		// 	getAllLandmark(c, db)
		// })
		routes.POST("/login", func(c *gin.Context) {
			UserLogin(c, db)
		})
	}
}

func UserLogin(c *gin.Context, db *gorm.DB) {
	// ใช้ DTO สำหรับรับ JSON
	var request dto.UserloginRequest

	// ตรวจสอบ JSON ว่าถูกต้องหรือไม่
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// ค้นหา customer
	var customer model.Customer
	if err := db.Where("email = ?", request.Email).First(&customer).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// ตรวจสอบรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(*request.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// ส่งข้อมูลที่อัปเดตกลับไป
	c.JSON(http.StatusOK, gin.H{
		"message":  "login success",
		"landmark": customer,
	})
}
