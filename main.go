package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	results   = make(map[int]int)
	resultsMu sync.Mutex
)

func main() {
	router := gin.Default()
	router.POST("/generate", generateNumbers)
	router.GET("/result/:id", getResult)

	// Start HTTP server in a separate goroutine
	go func() {
		if err := router.Run(":8080"); err != nil {
			fmt.Println("Failed to start HTTP server:", err)
		}
	}()

	// Graceful shutdown 
	gracefulShutdown()
}

func generateNumbers(c *gin.Context) {
	var requestBody struct {
		Amount int `json:"amount"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}


	numbers := generateRandomNumbers(requestBody.Amount)

	resultsMu.Lock()
	defer resultsMu.Unlock()
	results[len(results)] = sum(numbers)

	c.JSON(http.StatusOK, gin.H{"message": "Numbers generated successfully", "id": len(results) - 1})
}

func generateRandomNumbers(amount int) []int {
	numbers := make([]int, amount)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < amount; i++ {
		numbers[i] = rand.Intn(1000) 
	}
	return numbers
}

func sum(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

func getResult(c *gin.Context) {
	id := c.Param("id")
	resultID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	resultsMu.Lock()
	defer resultsMu.Unlock()
	result, ok := results[resultID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Result not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")
	os.Exit(0)
}