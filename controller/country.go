package controller

import (
	"final_go/dto"
	"final_go/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CountryController sets up country-related routes
func CountryController(router *gin.Engine, db *gorm.DB) {
	routes := router.Group("/countries")
	{
		routes.GET("/", func(c *gin.Context) {
			getAllCountry(c, db)
		})
		// routes.POST("/", createPerson)
	}
}

func getAllCountry(c *gin.Context, db *gorm.DB) {
	var countries []model.Country
	result := db.Find(&countries)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, countries)
}

func CreateCountry(c *gin.Context, db *gorm.DB) {
	var country dto.Country

	// รับข้อมูล JSON และเช็คว่าข้อมูลถูกต้องไหม
	if err := c.ShouldBindJSON(&country); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// บันทึกข้อมูลลงฐานข้อมูล
	result := db.Create(&country)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// ส่งข้อมูลกลับไป
	c.JSON(http.StatusOK, gin.H{
		"message": "Country created successfully",
	})
}
