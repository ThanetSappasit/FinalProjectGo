package main

import (
	"final_go/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	controller.StartServer()
}
