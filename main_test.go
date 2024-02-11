package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestGenerateNumbersEndpoint(t *testing.T) {
	router := setupRouter()


	requestBody := map[string]int{"amount": 5}
	requestBodyBytes, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/generate", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.Code)
	}

	var responseBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseBody)
	if _, ok := responseBody["id"]; !ok {
		t.Errorf("Expected response body to contain 'id'")
	}
}

func TestGetResultEndpoint(t *testing.T) {
	router := setupRouter()
	results[0] = 15 

	req, err := http.NewRequest("GET", "/result/0", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.Code)
	}

	var responseBody map[string]int
	json.NewDecoder(resp.Body).Decode(&responseBody)
	expectedResult := 15
	if responseBody["result"] != expectedResult {
		t.Errorf("Expected result %d, got %d", expectedResult, responseBody["result"])
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/generate", generateNumbers)
	router.GET("/result/:id", getResult)
	return router
}

func TestGenerateRandomNumbers(t *testing.T) {
	numbers := generateRandomNumbers(5)
	if len(numbers) != 5 {
		t.Errorf("Expected 5 numbers, got %d", len(numbers))
	}
}

func TestSum(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	expectedSum := 15
	actualSum := sum(numbers)
	if actualSum != expectedSum {
		t.Errorf("Expected sum %d, got %d", expectedSum, actualSum)
	}
}

func TestConcurrency(t *testing.T) {
	const numRequests = 100
	var wg sync.WaitGroup
	wg.Add(numRequests)

	
	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			requestBody := map[string]int{"amount": 5}
			requestBodyBytes, _ := json.Marshal(requestBody)
			req, _ := http.NewRequest("POST", "/generate", bytes.NewBuffer(requestBodyBytes))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router := setupRouter()
			router.ServeHTTP(resp, req)
		}()
	}

	wg.Wait()
	time.Sleep(1 * time.Second)


	resultsMu.Lock()
	defer resultsMu.Unlock()
	if len(results) -1 != numRequests {
		t.Errorf("Expected %d results, got %d", numRequests, len(results) - 1)
	}
}



func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}