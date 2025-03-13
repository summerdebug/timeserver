package main

import (
	"fmt"
	"net/http"
	"time"
)

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
	http.HandleFunc("/", timeHandler)
	fmt.Println("Server is running on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
}
