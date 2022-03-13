package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Order struct {
	Price  float64 `json:"price"`
	Volume uint64  `json:"volume"`
}

func main() {
	e := NewExchange()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/stock/:stock/*action", func(c *gin.Context) {
		stock := c.Param("stock")
		action := c.Param("action")
		if action == "add" {
			e.newStockType(stock)
		} else if action == "describe" {
			price, err := e.describeStock(stock)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": price})
		}
	})

	r.POST("/buy/:stock", func(c *gin.Context) {
		stock := c.Param("stock")
		var json Order
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err := e.buy(stock, json.Price, json.Volume)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "added buy order"})
	})

	r.POST("/sell/:stock", func(c *gin.Context) {
		stock := c.Param("stock")
		var json Order
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := e.addSell(stock, &Item{
			price:  json.Price,
			volume: json.Volume,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "added sell order"})
	})

	err := r.Run()
	if err != nil {
		panic(err)
	}
}
