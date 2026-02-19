package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	cache := NewCache()

	if err := cache.Refresh(); err != nil {
		log.Fatalf("Failed to load data from Mokky: %v", err)
	}

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			if err := cache.Refresh(); err != nil {
				log.Printf("Failed to refresh cache: %v", err)
			}
		}
	}()

	r := mux.NewRouter()
	r.Use(corsMiddleware)

	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
	})

	r.HandleFunc("/api/items", cache.ItemsHandler).Methods("GET")
	r.HandleFunc("/api/items/{uniqueId}", cache.ItemHandler).Methods("GET")
	r.HandleFunc("/api/items/{uniqueId}/similar", cache.SimilarHandler).Methods("GET")
	r.HandleFunc("/api/main-page/new", cache.MainPageNewHandler).Methods("GET")
	r.HandleFunc("/api/admin/items", cache.CreateItemHandler).Methods("POST")
	r.HandleFunc("/api/admin/items/{id}", cache.UpdateItemHandler).Methods("PUT", "PATCH")
	r.HandleFunc("/api/admin/items/{id}", cache.DeleteItemHandler).Methods("DELETE")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
