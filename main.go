package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("vim-go")
	router := gin.Default()

	//router.Static("/", "./views")
	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	api := router.Group("/api")
	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	api.GET("/ticker/:tickerID", TickerHandler)

	// Start and run the server
	router.Run(":3000")

}

func TickerHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"message": "Ticker handler not implemented yet",
		"note":    c.Param("tickerID"),
	})
}
