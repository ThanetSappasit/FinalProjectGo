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
func UserController(router *gin.Engine, db *gorm.DB) {
	routes := router.Group("/user")
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

func GetUserByEmail(c *gin.Context, db *gorm.DB) {
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

func UpdateUser(c *gin.Context, db *gorm.DB) {
	// ใช้ DTO สำหรับรับ JSON
	var request dto.UserUpdateRequest

	// ตรวจสอบ JSON ว่าถูกต้องหรือไม่
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// ค้นหา Landmark ที่ต้องการอัปเดต
	var customer model.Customer
	if err := db.Where("email = ?", request.Email).First(&customer).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// อัปเดตเฉพาะฟิลด์ที่ส่งมา
	updateData := map[string]interface{}{}

	if request.FirstName != nil {
		updateData["first_name"] = *request.FirstName
	}
	if request.LastName != nil {
		updateData["last_name"] = *request.LastName
	}
	if request.PhoneNumber != nil {
		updateData["phone_number"] = *request.PhoneNumber
	}
	if request.Address != nil {
		updateData["address"] = *request.Address
	}
	if request.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		updateData["password"] = string(hashedPassword)
	}

	// ถ้ามีค่าให้อัปเดต
	if len(updateData) > 0 {
		if err := db.Model(&customer).Updates(updateData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
	}

	// ส่งข้อมูลที่อัปเดตกลับไป
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    customer,
	})
}

func UpdatePassword(c *gin.Context, db *gorm.DB) {
	// ใช้ DTO สำหรับรับ JSON
	var request dto.UserChangePasswordRequest

	// ตรวจสอบ JSON ว่าถูกต้องหรือไม่
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// ตรวจสอบว่า email, old_password, new_password ไม่เป็น nil
	if request.Email == nil || request.Password == nil || request.NewPassword == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email, old password, and new password are required"})
		return
	}

	// ค้นหาผู้ใช้จากอีเมล
	var customer model.Customer
	if err := db.Where("email = ?", *request.Email).First(&customer).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
		return
	}

	// ตรวจสอบรหัสผ่านเก่าว่าตรงกับในฐานข้อมูลหรือไม่
	err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(*request.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect old password"})
		return
	}

	// เข้ารหัสรหัสผ่านใหม่
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	// อัปเดตรหัสผ่านใหม่ในฐานข้อมูล
	if err := db.Model(&customer).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// ส่งข้อความกลับไปว่าอัปเดตรหัสผ่านสำเร็จ
	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}
