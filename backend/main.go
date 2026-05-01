package main

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
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
				slog.Error("panic recovered", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func maxBytesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
		next.ServeHTTP(w, r)
	})
}

func generateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func csrfTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, err := r.Cookie("csrf_token")
			if err != nil {
				token := generateCSRFToken()
				http.SetCookie(w, &http.Cookie{
					Name:     "csrf_token",
					Value:    token,
					HttpOnly: false,
					Secure:   true,
					SameSite: http.SameSiteStrictMode,
					Path:     "/",
				})
			}
		}
		next.ServeHTTP(w, r)
	})
}

func csrfProtectionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}
		cookie, err := r.Cookie("csrf_token")
		if err != nil {
			http.Error(w, "CSRF token missing in cookie", http.StatusForbidden)
			return
		}
		headerToken := r.Header.Get("X-CSRF-Token")
		if headerToken == "" {
			http.Error(w, "CSRF token missing in header", http.StatusForbidden)
			return
		}
		if !hmac.Equal([]byte(cookie.Value), []byte(headerToken)) {
			http.Error(w, "CSRF token mismatch", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func newGlobalRateLimiter() func(http.Handler) http.Handler {
	rate := limiter.Rate{Period: 1 * time.Minute, Limit: 100}
	store := memory.NewStore()
	middleware := stdlib.NewMiddleware(limiter.New(store, rate))
	return middleware.Handler
}

func verifyRecaptcha(token string) (bool, error) {
	secret := os.Getenv("RECAPTCHA_SECRET")
	if secret == "" {
		return true, nil
	}
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{"secret": {secret}, "response": {token}})
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result struct {
		Success bool `json:"success"`
	}
	json.Unmarshal(body, &result)
	return result.Success, nil
}

func main() {
	initLogger()
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
				slog.Error("cache refresh failed", "error", err)
			}
		}
	}()

	rl := newRateLimiter(rate.Limit(5), 10)

	r := mux.NewRouter()
	r.Use(
		corsMiddleware,
		recoveryMiddleware,
		securityHeadersMiddleware,
		maxBytesMiddleware,
		csrfTokenMiddleware,
		csrfProtectionMiddleware,
		newGlobalRateLimiter(),
		authContextMiddleware,
		loggingMiddleware,
	)

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

	slog.Info("server started", "port", 8080)
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server started", "port", 8080)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server exited gracefully")
	}
}

func ExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
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
