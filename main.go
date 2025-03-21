package main

import (
	"context"
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

type TimeResponse struct {
	DayOfWeek  string `json:"day_of_week"`
	DayOfMonth int    `json:"day_of_month"`
	Month      string `json:"month"`
	Year       int    `json:"year"`
	Hour       int    `json:"hour"`
	Minute     int    `json:"minute"`
	Second     int    `json:"second"`
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
	t := time.Now()
	if r.Header.Get("Accept") == "application/json" {
		response := TimeResponse{
			DayOfWeek:  t.Weekday().String(),
			DayOfMonth: t.Day(),
			Month:      t.Month().String(),
			Year:       t.Year(),
			Hour:       t.Hour(),
			Minute:     t.Minute(),
			Second:     t.Second(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, t.Format(time.RFC3339))
	}
}

func TimeoutMiddleware(timeoutMs int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a context with timeout
			ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeoutMs)*time.Millisecond)
			defer cancel()

			// Attach the new context to the request
			r = r.WithContext(ctx)

			// Create a channel to detect when handler completes
			done := make(chan struct{})
			go func() {
				next.ServeHTTP(w, r)
				close(done)
			}()

			select {
			case <-ctx.Done():
				// Context timeout reached
				if ctx.Err() == context.DeadlineExceeded {
					http.Error(w, "Request timed out", http.StatusGatewayTimeout)
				}
			case <-done:
				// Handler finished within the timeout
			}
		})
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", timeHandler)

	timeoutMw := TimeoutMiddleware(1000)
	handlerChain := timeoutMw(loggingMiddleware(mux))

	err := http.ListenAndServe(":8080", handlerChain)
	if err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
}
