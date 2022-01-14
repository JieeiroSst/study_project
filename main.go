package main

import (
	"log"
	"net/http"
	_ "project/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Gin Swagger Example API
// @version 2.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http
func main() {
	// Gin instance
	r := gin.New()

	// Routes
	r.GET("/", HealthCheck)
	r.POST("/", HealthCheck1)

	url := ginSwagger.URL("http://localhost:3000/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Start server
	if err := r.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c *gin.Context) {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}

	c.JSON(http.StatusOK, res)
}

// HealthCheck1 godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} map[string]interface{}
// @Router / [post]
func HealthCheck1(c *gin.Context) {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}

	c.JSON(http.StatusOK, res)
}