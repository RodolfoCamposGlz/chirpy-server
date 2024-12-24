package main

import (
	"log"
	"net/http"
)

const port = "8080"
const filepathRoot = "."

func main (){
	mux := http.NewServeMux()

	// Create a new http.Server
	server := &http.Server{
		Addr:  ":" + port, // Bind to port 8080
		Handler: mux,     // Use the ServeMux as the handler
	}

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	mux.HandleFunc("/healthz",  func(w http.ResponseWriter, req *http.Request){
		w.Header().Set("Content-Type","text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	// Start the server
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}