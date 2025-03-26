package controller

import (
	"final_go/dto"
	"final_go/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CountryController sets up country-related routes
func LandmarkController(router *gin.Engine, db *gorm.DB) {
	routes := router.Group("/landmarks")
	{
		routes.GET("/", func(c *gin.Context) {
			getAllLandmark(c, db)
		})
		routes.POST("/", func(c *gin.Context) {
			CreateLandmark(c, db)
		})
		routes.DELETE("/", func(c *gin.Context) {
			DeleteLandmark(c, db)
		})
		routes.PUT("/", func(c *gin.Context) {
			UpdateLandmark(c, db)
		})
	}
}

func getAllLandmark(c *gin.Context, db *gorm.DB) {
	var landmarks []model.Landmark
	result := db.Joins("CountryData").Find(&landmarks)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, landmarks)
}

func CreateLandmark(c *gin.Context, db *gorm.DB) {
	var landmark dto.Landmark // ใช้ dto.Landmark โดยตรง

	// รับ JSON และตรวจสอบว่าข้อมูลถูกต้องหรือไม่
	if err := c.ShouldBindJSON(&landmark); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ตรวจสอบว่า Country มีอยู่จริงหรือไม่
	var country model.Country
	if err := db.First(&country, landmark.Country).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Country ID"})
		return
	}

	// ใช้ dto.Landmark บันทึกลงฐานข้อมูลโดยตรง
	if err := db.Table("landmark").Create(&landmark).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ส่งข้อมูลกลับไป
	c.JSON(http.StatusCreated, gin.H{
		"message": "Landmark created successfully",
	})
}
func DeleteLandmark(c *gin.Context, db *gorm.DB) {
	// สร้าง struct สำหรับรับ JSON
	var request struct {
		Idx int `json:"idx"`
	}

	// อ่าน JSON และตรวจสอบความถูกต้อง
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// ตรวจสอบว่ามี Landmark ที่ต้องการลบหรือไม่
	var landmark model.Landmark
	if err := db.First(&landmark, request.Idx).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Landmark not found"})
		return
	}

	// ลบข้อมูล
	if err := db.Delete(&landmark).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete landmark"})
		return
	}

	// ตอบกลับว่า ลบสำเร็จ
	c.JSON(http.StatusOK, gin.H{"message": "Landmark deleted successfully"})
}

func UpdateLandmark(c *gin.Context, db *gorm.DB) {
	// ใช้ DTO สำหรับรับ JSON
	var request dto.UpdateLandmarkRequest

	// ตรวจสอบ JSON ว่าถูกต้องหรือไม่
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// ค้นหา Landmark ที่ต้องการอัปเดต
	var landmark model.Landmark
	if err := db.First(&landmark, request.Idx).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Landmark not found"})
		return
	}

	// อัปเดตเฉพาะฟิลด์ที่ส่งมา
	updateData := map[string]interface{}{}

	if request.Name != nil {
		updateData["name"] = *request.Name
	}
	if request.Country != nil {
		updateData["country"] = *request.Country
	}
	if request.Detail != nil {
		updateData["detail"] = *request.Detail
	}
	if request.Url != nil {
		updateData["url"] = *request.Url
	}

	// ถ้ามีค่าให้อัปเดต
	if len(updateData) > 0 {
		if err := db.Model(&landmark).Updates(updateData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update landmark"})
			return
		}
	}

	// ส่งข้อมูลที่อัปเดตกลับไป
	c.JSON(http.StatusOK, gin.H{
		"message":  "Landmark updated successfully",
		"landmark": landmark,
	})
}
