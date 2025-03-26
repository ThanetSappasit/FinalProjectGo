package controller

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connection() (*gorm.DB, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	dsn := viper.GetString("mysql.dsn")
	if dsn == "" {
		return nil, fmt.Errorf("mysql.dsn is empty")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	fmt.Println("Database connection successful")
	return db, nil
}

func StartServer() {
	router := gin.Default()

	// Establish database connection
	db, err := Connection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "API is now working"})
	})

	// CountryController(router, db)  // Pass the actual database connection
	// LandmarkController(router, db) // Pass the actual database connection
	LoginController(router, db)
	UserController(router, db)
	ProductController(router, db)
	CartController(router, db)

	// Start the server
	router.Run()
}
