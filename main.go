package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	requestCount    int
	status5xxCount  int
	status200Count  int
	countMutex      sync.Mutex
	serverStartTime time.Time
	oneMinute       = time.Minute
)

type CountResponse struct {
	TotalRequests int `json:"total_requests"`
	Status5xx     int `json:"status_5xx"`
	Status200     int `json:"status_200"`
}

type Response struct {
	CurrentTime string `json:"current_time"`
	Message     string `json:"message"`
	IPAddress   string `json:"ip_address"`
}

func hiHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	ipAddress := r.RemoteAddr

	response := Response{
		CurrentTime: currentTime,
		Message:     "Hello!",
		IPAddress:   ipAddress,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	countMutex.Lock()
	defer countMutex.Unlock()

	currentTime := time.Now()

	requestCount++
	if currentTime.Sub(serverStartTime) <= oneMinute {
		status5xxCount++
		fmt.Printf("Service Unavailable")
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	status200Count++
	fmt.Printf("Service OK")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func countHandler(w http.ResponseWriter, r *http.Request) {
	countMutex.Lock()
	defer countMutex.Unlock()

	response := CountResponse{
		TotalRequests: requestCount,
		Status5xx:     status5xxCount,
		Status200:     status200Count,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleSIGTERM() {
	fmt.Println("Received SIGTERM signal. Graceful shutdown...")
	// 在这里添加你的清理逻辑，例如关闭连接、保存状态等
	time.Sleep(10 * time.Second)
	os.Exit(0)
}

func main() {
	// 注册 SIGTERM 信号处理函数
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalCh
		handleSIGTERM()
	}()

	http.HandleFunc("/hi", hiHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/count", countHandler)
	serverStartTime = time.Now()
	port := "8081"
	fmt.Printf("Server listening on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
