package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
		if allowedOrigin == "" {
			allowedOrigin = "http://localhost:5173" // для разработки
		}
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := InitDB(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer DB.Close()

	cache := NewDBCache()
	if err := cache.Refresh(); err != nil {
		log.Fatalf("Failed to load data from DB: %v", err)
	}

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			if err := cache.Refresh(); err != nil {
				log.Printf("Failed to refresh cache: %v", err)
			}
		}
	}()

	rl := newRateLimiter(rate.Limit(5), 10)

	r := mux.NewRouter()
	r.Use(corsMiddleware, recoveryMiddleware, securityHeadersMiddleware)

	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	r.HandleFunc("/api/items", cache.ItemsHandler).Methods("GET")
	r.HandleFunc("/api/items/{uniqueId}", cache.ItemHandler).Methods("GET")
	r.HandleFunc("/api/items/{uniqueId}/similar", cache.SimilarHandler).Methods("GET")
	r.HandleFunc("/api/main-page/new", cache.MainPageNewHandler).Methods("GET")
	r.HandleFunc("/api/filters", cache.FiltersHandler).Methods("GET")
	r.HandleFunc("/api/exchange-rate", ExchangeRateHandler).Methods("GET")

	r.HandleFunc("/api/auth/login", rl.middleware(LoginHandler, rate.Limit(2), 5)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/register", rl.middleware(RegisterHandler, rate.Limit(2), 5)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/me", MeHandler).Methods("GET", "OPTIONS")

	r.HandleFunc("/api/user/profile", UpdateProfileHandler).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/user/profile", DeleteAccountHandler).Methods("DELETE", "OPTIONS")

	r.HandleFunc("/api/user/favourites", cache.GetUserFavourites).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/user/favourites/{uniqueId}", cache.AddToFavourites).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/user/favourites/{uniqueId}", cache.RemoveFromFavourites).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/api/user/favourites/items", cache.GetFavouriteItems).Methods("GET", "OPTIONS")

	r.HandleFunc("/api/user/cart", cache.GetUserCart).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/user/cart/{uniqueId}", cache.AddToCart).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/user/cart/{uniqueId}", cache.UpdateCartItem).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/user/cart/{uniqueId}", cache.RemoveFromCart).Methods("DELETE", "OPTIONS")

	r.HandleFunc("/api/user/orders", cache.GetUserOrders).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/user/orders", cache.CreateOrderHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/user/orders/{id}", cache.GetOrderDetails).Methods("GET", "OPTIONS")

	r.HandleFunc("/api/admin/items", adminOnly(cache.AdminItemsHandler)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/admin/items/{id}", adminOnly(cache.AdminItemHandler)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/admin/items", adminOnly(cache.CreateItemHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/admin/items/{id}", adminOnly(cache.UpdateItemHandler)).Methods("PUT", "PATCH", "OPTIONS")
	r.HandleFunc("/api/admin/items/{id}", adminOnly(cache.DeleteItemHandler)).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/api/admin/orders", adminOnly(cache.AdminGetOrders)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/admin/orders/{id}/items", adminOnly(cache.AdminGetOrderItems)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/admin/order-items/{itemId}/status", adminOnly(cache.AdminUpdateOrderItemStatus)).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/admin/orders/{id}", adminOnly(cache.AdminGetOrderDetails)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/admin/orders/{id}/status", adminOnly(cache.AdminUpdateOrderStatus)).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/admin/orders/{id}/status/last", adminOnly(cache.AdminDeleteLastOrderStatus)).Methods("DELETE", "OPTIONS")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func ExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	// Заглушка
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"rate":11.5}`))
}

func adminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		isAdmin, err := isUserAdmin(userID)
		if err != nil || !isAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
