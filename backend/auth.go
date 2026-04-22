package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("WARNING: JWT_SECRET not set, using random secret. All tokens will be invalid after server restart.")
		secretBytes := make([]byte, 32)
		_, err := rand.Read(secretBytes)
		if err != nil {
			panic("failed to generate random JWT secret")
		}
		jwtSecret = secretBytes
	} else {
		jwtSecret = []byte(secret)
	}

	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(logFile)
	}
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func generateToken(userID int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func getUserIDFromToken(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("missing Authorization header")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return 0, fmt.Errorf("Bearer prefix missing")
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}
	return 0, fmt.Errorf("invalid token")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if len(req.Email) > 255 || !strings.Contains(req.Email, "@") {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}
	if len(req.Name) > 100 {
		http.Error(w, "Name too long", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	var userID int
	err = DB.QueryRow(context.Background(), `
        INSERT INTO users (email, password, name, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id
    `, req.Email, string(hashed), req.Name).Scan(&userID)
	if err != nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}
	token, _ := generateToken(userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       userID,
			"email":    req.Email,
			"name":     req.Name,
			"is_admin": false,
		},
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var userID int
	var hashedPassword string
	var name string
	var isAdmin bool
	err := DB.QueryRow(context.Background(),
		"SELECT id, password, name, is_admin FROM users WHERE email = $1", req.Email).
		Scan(&userID, &hashedPassword, &name, &isAdmin)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	token, _ := generateToken(userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       userID,
			"email":    req.Email,
			"name":     name,
			"is_admin": isAdmin,
		},
	})
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var email, name string
	var isAdmin bool
	err = DB.QueryRow(context.Background(), "SELECT email, name, is_admin FROM users WHERE id = $1", userID).
		Scan(&email, &name, &isAdmin)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       userID,
		"email":    email,
		"name":     name,
		"is_admin": isAdmin,
	})
}

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Email != "" {
		var exists bool
		err := DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id != $2)", req.Email, userID).Scan(&exists)
		if err == nil && exists {
			http.Error(w, "Email already taken", http.StatusConflict)
			return
		}
	}

	var updates []string
	var args []interface{}
	argPos := 1
	if req.Name != "" {
		updates = append(updates, fmt.Sprintf("name = $%d", argPos))
		args = append(args, req.Name)
		argPos++
	}
	if req.Email != "" {
		updates = append(updates, fmt.Sprintf("email = $%d", argPos))
		args = append(args, req.Email)
		argPos++
	}
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		updates = append(updates, fmt.Sprintf("password = $%d", argPos))
		args = append(args, string(hashed))
		argPos++
	}
	if len(updates) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(updates, ", "), argPos)
	args = append(args, userID)
	_, err = DB.Exec(context.Background(), query, args...)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	var email, name string
	DB.QueryRow(context.Background(), "SELECT email, name FROM users WHERE id = $1", userID).Scan(&email, &name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    userID,
		"email": email,
		"name":  name,
	})
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	_, err = DB.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func isUserAdmin(userID int) (bool, error) {
	var isAdmin bool
	err := DB.QueryRow(context.Background(), "SELECT is_admin FROM users WHERE id = $1", userID).Scan(&isAdmin)
	return isAdmin, err
}
