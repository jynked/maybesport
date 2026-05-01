package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (c *DBCache) GetUserCart(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	rows, err := DB.Query(context.Background(), `
    	SELECT pi.unique_id, pi.images, pi.color_ru, pi.color_en,
           p.title_ru, p.title_en, p.id,
           uc.size, uc.quantity,
           COALESCE(pis.price, 0) as price,
           COALESCE(pis.is_on_request, false) as is_on_request,
           COALESCE(pis.quantity, 0) as stock
		FROM user_cart uc
		JOIN product_items pi ON uc.product_item_id = pi.id
		JOIN products p ON pi.product_id = p.id
		LEFT JOIN product_item_sizes pis ON pis.product_item_id = pi.id AND pis.size = uc.size
		WHERE uc.user_id = $1
	`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result []CartItemResponse
	for rows.Next() {
		var cr CartItemResponse
		var images []string
		var titleRu, titleEn string
		var stock int64
		err := rows.Scan(&cr.UniqueId, &images, &cr.Color.Ru, &cr.Color.En,
			&titleRu, &titleEn, &cr.ID, &cr.Size, &cr.Quantity,
			&cr.Price, &cr.IsOnRequest, &stock)
		if err != nil {
			continue
		}
		cr.Title = Lang{Ru: titleRu, En: titleEn}
		if len(images) > 0 {
			cr.Image = images[0]
		}
		cr.Stock = stock
		cr.Availability = "out_of_stock"
		if stock > 0 && !cr.IsOnRequest {
			cr.Availability = "available"
		} else if stock > 0 && cr.IsOnRequest {
			cr.Availability = "on_request"
		}
		result = append(result, cr)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (c *DBCache) AddToCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uniqueId := vars["uniqueId"]
	if uniqueId == "" {
		http.Error(w, "uniqueId required", http.StatusBadRequest)
		return
	}
	var body struct {
		Size     interface{} `json:"size"`
		Quantity int64       `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if body.Size == nil || body.Quantity < 1 {
		http.Error(w, "size and positive quantity required", http.StatusBadRequest)
		return
	}
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	sizeStr := toString(body.Size)

	var productItemID int
	var currentStock int64
	var isOnRequest bool
	err = DB.QueryRow(context.Background(), `
        SELECT pi.id, pis.quantity, pis.is_on_request
        FROM product_items pi
        JOIN product_item_sizes pis ON pis.product_item_id = pi.id
        WHERE pi.unique_id = $1 AND pis.size = $2
    `, uniqueId, sizeStr).Scan(&productItemID, &currentStock, &isOnRequest)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if body.Quantity > currentStock {
		http.Error(w, "Requested quantity exceeds available stock", http.StatusBadRequest)
		return
	}

	if !isOnRequest && body.Quantity > currentStock {
		http.Error(w, fmt.Sprintf("Not enough stock. Available: %d", currentStock), http.StatusConflict)
		return
	}

	_, err = DB.Exec(context.Background(), `
        INSERT INTO user_cart (user_id, product_item_id, size, quantity, updated_at)
        VALUES ($1, $2, $3, $4, NOW())
        ON CONFLICT (user_id, product_item_id, size) DO UPDATE SET quantity = EXCLUDED.quantity, updated_at = NOW()
    `, userID, productItemID, sizeStr, body.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *DBCache) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uniqueId := vars["uniqueId"]
	if uniqueId == "" {
		http.Error(w, "uniqueId required", http.StatusBadRequest)
		return
	}
	var body struct {
		Size     interface{} `json:"size"`
		Quantity int64       `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if body.Size == nil {
		http.Error(w, "size required", http.StatusBadRequest)
		return
	}
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	sizeStr := toString(body.Size)

	var productItemID int
	var currentStock int64
	var isOnRequest bool
	err = DB.QueryRow(context.Background(), `
		SELECT pi.id, pis.quantity, pis.is_on_request
		FROM product_items pi
		JOIN product_item_sizes pis ON pis.product_item_id = pi.id
		WHERE pi.unique_id = $1 AND pis.size = $2
	`, uniqueId, sizeStr).Scan(&productItemID, &currentStock, &isOnRequest)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if body.Quantity <= 0 {
		_, err = DB.Exec(context.Background(), `
			DELETE FROM user_cart WHERE user_id=$1 AND product_item_id=$2 AND size=$3
		`, userID, productItemID, sizeStr)
	} else {
		_, err = DB.Exec(context.Background(), `
			UPDATE user_cart SET quantity=$4, updated_at=NOW()
			WHERE user_id=$1 AND product_item_id=$2 AND size=$3
		`, userID, productItemID, sizeStr, body.Quantity)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *DBCache) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uniqueId := vars["uniqueId"]
	sizeParam := r.URL.Query().Get("size")
	if uniqueId == "" || sizeParam == "" {
		http.Error(w, "uniqueId and size required", http.StatusBadRequest)
		return
	}
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var productItemID int
	err = DB.QueryRow(context.Background(), "SELECT id FROM product_items WHERE unique_id = $1", uniqueId).Scan(&productItemID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	_, err = DB.Exec(context.Background(), `
		DELETE FROM user_cart WHERE user_id=$1 AND product_item_id=$2 AND size=$3
	`, userID, productItemID, sizeParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
