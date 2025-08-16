package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var store = make(map[string]string)
var mu sync.RWMutex

func putHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	// Error check if key or value is empty
	if key == "" || value == "" {
		http.Error(w, "Key and value are required", http.StatusBadRequest)
		return
	}

	// lock the map
	mu.Lock()

	// Put the key and value in the map
	store[key] = value
	mu.Unlock()

	fmt.Fprintf(w, "Stored %s -> %s", key, value)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	// error check
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}

	// lock the map
	mu.RLock()

	// Get the value from the map
	value, ok := store[key]
	mu.RUnlock()

	if !ok {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{key: value})
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	// error check
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}

	// lock the map
	mu.Lock()

	// Check if the key exists in the map
	_, ok := store[key]
	if ok {
		// Delete the key from the map
		delete(store, key)
	} else {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	mu.Unlock()

	fmt.Fprintf(w, "Deleted %s", key)
}

func main() {
	// Register the handlers for the endpoints
	http.HandleFunc("/put", putHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/delete", deleteHandler)

	// Get the port from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	// Start the server
	fmt.Printf("Starting server on port %s\n", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
