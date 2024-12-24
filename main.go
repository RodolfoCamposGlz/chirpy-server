package main

import (
	"fmt"
	"net/http"
)


func main (){
	mux := http.NewServeMux()

	// Create a new http.Server
	server := &http.Server{
		Addr:    ":8080", // Bind to port 8080
		Handler: mux,     // Use the ServeMux as the handler
	}

	// Start the server
	fmt.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
	}
}