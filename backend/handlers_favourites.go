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
        SELECT pi.unique_id
        FROM user_favourites uf
        JOIN product_items pi ON uf.product_item_id = pi.id
        WHERE uf.user_id = $1
    `, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var favs []FavouriteItem
	for rows.Next() {
		var fav FavouriteItem
		rows.Scan(&fav.UniqueId)
		favs = append(favs, fav)
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
        INSERT INTO user_favourites (user_id, product_item_id, added_at)
        VALUES ($1, $2, NOW())
        ON CONFLICT (user_id, product_item_id) DO NOTHING
    `, userID, productItemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *DBCache) RemoveFromFavourites(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uniqueId := vars["uniqueId"]
	if uniqueId == "" {
		http.Error(w, "uniqueId required", http.StatusBadRequest)
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
        DELETE FROM user_favourites WHERE user_id=$1 AND product_item_id=$2
    `, userID, productItemID)
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
		SELECT pi.unique_id, pi.images, p.id, pi.title_ru, pi.title_en,
			MIN(pis.price_cny) as min_price_cny,
			BOOL_OR(pis.is_on_request) as has_on_request,
			MIN(CASE WHEN pis.quantity > 0 AND NOT pis.is_on_request THEN pis.price_cny ELSE NULL END) as available_price_cny
		FROM user_favourites uf
		JOIN product_items pi ON uf.product_item_id = pi.id
		JOIN products p ON pi.product_id = p.id
		LEFT JOIN product_item_sizes pis ON pis.product_item_id = pi.id
		WHERE uf.user_id = $1
		GROUP BY pi.unique_id, pi.images, p.id, pi.title_ru, pi.title_en
		ORDER BY MAX(uf.added_at) DESC
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
		var minPriceCny int64
		var hasOnRequest bool
		var availablePriceCny *int64
		err := rows.Scan(&fr.UniqueId, &images, &fr.ID, &titleRu, &titleEn, &minPriceCny, &hasOnRequest, &availablePriceCny)
		if err != nil {
			continue
		}
		fr.Title = Lang{Ru: titleRu, En: titleEn}
		if len(images) > 0 {
			fr.Image = images[0]
		}
		fr.Price = ConvertCnyToRub(minPriceCny)
		fr.IsOnRequest = hasOnRequest
		fr.Quantity = 0
		fr.Availability = "out_of_stock"
		if availablePriceCny != nil {
			fr.Availability = "available"
		} else if hasOnRequest {
			fr.Availability = "on_request"
		}
		result = append(result, fr)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
