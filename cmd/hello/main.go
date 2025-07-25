package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Message   string    `json:"message"`
}

type HelloResponse struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Message:   "Hello World API is running",
	}
	
	json.NewEncoder(w).Encode(response)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	response := HelloResponse{
		Message:   "Hello, World from Kubernetes!",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}
	
	json.NewEncoder(w).Encode(response)
}

func apiHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Message:   "API v1 is healthy",
	}
	
	json.NewEncoder(w).Encode(response)
}

func main() {
	// 设置路由
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/v1/health", apiHealthHandler)
	
	// 启动服务器
	port := "8080"
	fmt.Printf("Hello World API server starting on port %s...\n", port)
	fmt.Printf("Health check: http://localhost:%s/health\n", port)
	fmt.Printf("API health: http://localhost:%s/api/v1/health\n", port)
	fmt.Printf("Hello endpoint: http://localhost:%s/\n", port)
	
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
