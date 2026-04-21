package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("CORS middleware: %s %s", r.Method, r.URL.Path)
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
	r.Use(corsMiddleware, recoveryMiddleware)

	initExchangeRate()
	startExchangeRateUpdater()

	r.HandleFunc("/api/exchange-rate", cache.ExchangeRateHandler).Methods("GET")

	r.HandleFunc("/api/items", cache.ItemsHandler).Methods("GET")
	r.HandleFunc("/api/items/{uniqueId}", cache.ItemHandler).Methods("GET")
	r.HandleFunc("/api/items/{uniqueId}/similar", cache.SimilarHandler).Methods("GET")

	r.HandleFunc("/api/main-page/new", cache.MainPageNewHandler).Methods("GET")

	r.HandleFunc("/api/admin/items", cache.CreateItemHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/admin/items", cache.AdminItemsHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/admin/items/{id}", cache.AdminItemHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/admin/items/{id}", cache.UpdateItemHandler).Methods("PUT", "PATCH", "OPTIONS")
	r.HandleFunc("/api/admin/items/{id}", cache.DeleteItemHandler).Methods("DELETE", "OPTIONS")

	r.HandleFunc("/api/auth/login", LoginHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/register", RegisterHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/me", MeHandler).Methods("GET", "OPTIONS")

	r.HandleFunc("/api/filters", cache.FiltersHandler).Methods("GET")

	r.HandleFunc("/api/user/favourites", cache.GetUserFavourites).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/user/favourites/{uniqueId}", cache.AddToFavourites).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/user/favourites/{uniqueId}", cache.RemoveFromFavourites).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/api/user/favourites/items", cache.GetFavouriteItems).Methods("GET", "OPTIONS")

	r.HandleFunc("/api/user/orders", cache.GetUserOrders).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/user/orders/{id}", cache.GetOrderDetails).Methods("GET", "OPTIONS")

	r.HandleFunc("/api/user/cart", cache.GetUserCart).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/user/cart/{uniqueId}", cache.AddToCart).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/user/cart/{uniqueId}", cache.UpdateCartItem).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/user/cart/{uniqueId}", cache.RemoveFromCart).Methods("DELETE", "OPTIONS")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
