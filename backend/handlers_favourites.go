package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (c *DBCache) GetUserFavourites(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	rows, err := DB.Query(context.Background(), `
		SELECT product_item_id, size FROM user_favourites WHERE user_id = $1
	`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var favs []FavouriteItem
	for rows.Next() {
		var itemID int
		var size string
		if err := rows.Scan(&itemID, &size); err != nil {
			continue
		}
		var uniqueID string
		DB.QueryRow(context.Background(), "SELECT unique_id FROM product_items WHERE id = $1", itemID).Scan(&uniqueID)
		favs = append(favs, FavouriteItem{UniqueId: uniqueID, Size: size})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(favs)
}

func (c *DBCache) AddToFavourites(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uniqueId := vars["uniqueId"]
	if uniqueId == "" {
		http.Error(w, "uniqueId required", http.StatusBadRequest)
		return
	}
	var body struct {
		Size interface{} `json:"size"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Size == nil {
		http.Error(w, "size required", http.StatusBadRequest)
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
	sizeStr := toString(body.Size)
	_, err = DB.Exec(context.Background(), `
		INSERT INTO user_favourites (user_id, product_item_id, size, added_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, product_item_id, size) DO NOTHING
	`, userID, productItemID, sizeStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *DBCache) RemoveFromFavourites(w http.ResponseWriter, r *http.Request) {
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
		DELETE FROM user_favourites WHERE user_id=$1 AND product_item_id=$2 AND size=$3
	`, userID, productItemID, sizeParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *DBCache) GetFavouriteItems(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	rows, err := DB.Query(context.Background(), `
		SELECT pi.unique_id, pi.images, p.id, p.title_ru, p.title_en, 
		       uf.size, pis.price, pis.is_on_request, pis.quantity
		FROM user_favourites uf
		JOIN product_items pi ON uf.product_item_id = pi.id
		JOIN products p ON pi.product_id = p.id
		LEFT JOIN product_item_sizes pis ON pis.product_item_id = pi.id AND pis.size = uf.size
		WHERE uf.user_id = $1
	`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var result []FavouriteItemResponse
	for rows.Next() {
		var fr FavouriteItemResponse
		var images []string
		var titleRu, titleEn string
		err := rows.Scan(&fr.UniqueId, &images, &fr.ID, &titleRu, &titleEn, &fr.Size, &fr.Price, &fr.IsOnRequest, &fr.Quantity)
		if err != nil {
			continue
		}
		fr.Title = Lang{Ru: titleRu, En: titleEn}
		if len(images) > 0 {
			fr.Image = images[0]
		}
		fr.Availability = "out_of_stock"
		if fr.Quantity > 0 && !fr.IsOnRequest {
			fr.Availability = "available"
		} else if fr.Quantity > 0 && fr.IsOnRequest {
			fr.Availability = "on_request"
		}
		result = append(result, fr)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
