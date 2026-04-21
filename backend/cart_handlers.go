package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (c *Cache) GetUserCart(w http.ResponseWriter, r *http.Request) {
	user, err := getAuthenticatedUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	authHeader := r.Header.Get("Authorization")
	fullUser, err := fetchUserFromMokky(user.ID, authHeader)
	if err != nil {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	if len(fullUser.Cart) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]CartItemResponse{})
		return
	}

	allItems := c.GetFlatItems()
	result := []CartItemResponse{}

	for _, cartItem := range fullUser.Cart {
		var targetItem *ItemFlatten
		for _, it := range allItems {
			if it.UniqueId == cartItem.UniqueId {
				targetItem = &it
				break
			}
		}
		if targetItem == nil {
			continue
		}

		var selectedSize Size
		found := false
		for _, sz := range targetItem.Sizes {
			if toString(sz.Size) == toString(cartItem.Size) {
				selectedSize = sz
				found = true
				break
			}
		}
		if !found {
			continue
		}

		image := ""
		if len(targetItem.Images) > 0 {
			image = targetItem.Images[0]
		}

		availability := "out_of_stock"
		if selectedSize.Quantity > 0 && !selectedSize.IsOnRequest {
			availability = "available"
		} else if selectedSize.Quantity > 0 && selectedSize.IsOnRequest {
			availability = "on_request"
		}

		resp := CartItemResponse{
			UniqueId:     targetItem.UniqueId,
			ID:           targetItem.ID,
			Title:        targetItem.Title,
			Image:        image,
			Size:         selectedSize.Size,
			Price:        selectedSize.Price,
			Quantity:     cartItem.Quantity,
			IsOnRequest:  selectedSize.IsOnRequest,
			Stock:        selectedSize.Quantity,
			Availability: availability,
		}
		result = append(result, resp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (c *Cache) AddToCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uniqueId := vars["uniqueId"]
	if uniqueId == "" {
		http.Error(w, "uniqueId required", http.StatusBadRequest)
		return
	}

	var body struct {
		Size     interface{} `json:"size"`
		Quantity int         `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if body.Size == nil || body.Quantity < 1 {
		http.Error(w, "size and positive quantity required", http.StatusBadRequest)
		return
	}

	user, err := getAuthenticatedUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	authHeader := r.Header.Get("Authorization")
	fullUser, err := fetchUserFromMokky(user.ID, authHeader)
	if err != nil {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	if fullUser.Cart == nil {
		fullUser.Cart = []CartItem{}
	}

	found := false
	for i, item := range fullUser.Cart {
		if item.UniqueId == uniqueId && toString(item.Size) == toString(body.Size) {
			fullUser.Cart[i].Quantity += body.Quantity
			found = true
			break
		}
	}
	if !found {
		fullUser.Cart = append(fullUser.Cart, CartItem{
			UniqueId: uniqueId,
			Size:     body.Size,
			Quantity: body.Quantity,
		})
	}

	if err := updateUserInMokky(fullUser, authHeader); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Cache) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uniqueId := vars["uniqueId"]
	if uniqueId == "" {
		http.Error(w, "uniqueId required", http.StatusBadRequest)
		return
	}

	var body struct {
		Size     interface{} `json:"size"`
		Quantity int         `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if body.Size == nil {
		http.Error(w, "size required", http.StatusBadRequest)
		return
	}

	user, err := getAuthenticatedUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	authHeader := r.Header.Get("Authorization")
	fullUser, err := fetchUserFromMokky(user.ID, authHeader)
	if err != nil {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	newCart := []CartItem{}
	for _, item := range fullUser.Cart {
		if item.UniqueId == uniqueId && toString(item.Size) == toString(body.Size) {
			if body.Quantity > 0 {
				item.Quantity = body.Quantity
				newCart = append(newCart, item)
			}
			continue
		}
		newCart = append(newCart, item)
	}
	fullUser.Cart = newCart

	if err := updateUserInMokky(fullUser, authHeader); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Cache) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uniqueId := vars["uniqueId"]
	sizeParam := r.URL.Query().Get("size")
	if uniqueId == "" || sizeParam == "" {
		http.Error(w, "uniqueId and size are required", http.StatusBadRequest)
		return
	}

	user, err := getAuthenticatedUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	authHeader := r.Header.Get("Authorization")
	fullUser, err := fetchUserFromMokky(user.ID, authHeader)
	if err != nil {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	newCart := []CartItem{}
	for _, item := range fullUser.Cart {
		if !(item.UniqueId == uniqueId && toString(item.Size) == sizeParam) {
			newCart = append(newCart, item)
		}
	}
	fullUser.Cart = newCart

	if err := updateUserInMokky(fullUser, authHeader); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type CartItemResponse struct {
	UniqueId     string      `json:"uniqueId"`
	ID           int         `json:"id"`
	Title        Lang        `json:"title"`
	Image        string      `json:"image"`
	Size         interface{} `json:"size"`
	Price        int         `json:"price"`
	Quantity     int         `json:"quantity"`
	IsOnRequest  bool        `json:"isOnRequest"`
	Stock        int         `json:"stock"`
	Availability string      `json:"availability"`
}
