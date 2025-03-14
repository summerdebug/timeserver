package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	IP        string `json:"ip"`
	Method    string `json:"method"`
	Path      string `json:"path"`
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logEntry := LogEntry{
			Timestamp: time.Now().Format(time.RFC3339),
			IP:        r.RemoteAddr,
			Method:    r.Method,
			Path:      r.URL.Path,
		}
		logData, _ := json.Marshal(logEntry)
		fmt.Println(string(logData))
		next.ServeHTTP(w, r)
	})
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	timeNow := time.Now().Format(time.RFC3339)
	_, err := fmt.Fprintln(w, timeNow)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", timeHandler)

	loggedMux := loggingMiddleware(mux)

	err := http.ListenAndServe(":8080", loggedMux)
	if err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
}
