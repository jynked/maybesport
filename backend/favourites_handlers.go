package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (c *Cache) GetUserFavourites(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fullUser.Favourites)
}

func (c *Cache) AddToFavourites(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uniqueId := vars["uniqueId"]
	if uniqueId == "" {
		http.Error(w, "uniqueId required", http.StatusBadRequest)
		return
	}

	var body struct {
		Size interface{} `json:"size"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body: size required", http.StatusBadRequest)
		return
	}
	if body.Size == nil {
		http.Error(w, "size is required", http.StatusBadRequest)
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

	if fullUser.Favourites == nil {
		fullUser.Favourites = []FavouriteItem{}
	}

	already := false
	for _, fav := range fullUser.Favourites {
		if fav.UniqueId == uniqueId && toString(fav.Size) == toString(body.Size) {
			already = true
			break
		}
	}
	if !already {
		fullUser.Favourites = append(fullUser.Favourites, FavouriteItem{
			UniqueId: uniqueId,
			Size:     body.Size,
		})
		if err := updateUserInMokky(fullUser, authHeader); err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Cache) RemoveFromFavourites(w http.ResponseWriter, r *http.Request) {
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

	newFavs := []FavouriteItem{}
	for _, fav := range fullUser.Favourites {
		if fav.UniqueId != uniqueId || toString(fav.Size) != sizeParam {
			newFavs = append(newFavs, fav)
		}
	}
	fullUser.Favourites = newFavs
	if err := updateUserInMokky(fullUser, authHeader); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Cache) GetFavouriteItems(w http.ResponseWriter, r *http.Request) {
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

	if len(fullUser.Favourites) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]FavouriteItemResponse{})
		return
	}

	allItems := c.GetFlatItems()
	result := []FavouriteItemResponse{}

	for _, fav := range fullUser.Favourites {
		var targetItem *ItemFlatten
		for _, it := range allItems {
			if it.UniqueId == fav.UniqueId {
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
			if toString(sz.Size) == toString(fav.Size) {
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

		resp := FavouriteItemResponse{
			UniqueId:     targetItem.UniqueId,
			ID:           targetItem.ID,
			Title:        targetItem.Title,
			Image:        image,
			Size:         selectedSize.Size,
			Price:        selectedSize.Price,
			IsOnRequest:  selectedSize.IsOnRequest,
			Quantity:     selectedSize.Quantity,
			Availability: availability,
		}
		result = append(result, resp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getAuthenticatedUser(r *http.Request) (*User, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("no auth header")
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://3b7b2b24dfd8c527.mokky.dev/auth_me", nil)
	req.Header.Set("Authorization", authHeader)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unauthorized")
	}
	defer resp.Body.Close()
	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func fetchUserFromMokky(userID int, authHeader string) (*User, error) {
	url := fmt.Sprintf("https://3b7b2b24dfd8c527.mokky.dev/users/%d", userID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", authHeader)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user not found, status %d", resp.StatusCode)
	}

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	user := &User{}
	if id, ok := raw["id"].(float64); ok {
		user.ID = int(id)
	}
	if email, ok := raw["email"].(string); ok {
		user.Email = email
	}
	if pass, ok := raw["password"].(string); ok {
		user.Password = pass
	}

	if favRaw, ok := raw["favourites"]; ok && favRaw != nil {
		if arr, ok := favRaw.([]interface{}); ok {
			newFavs := []FavouriteItem{}
			for _, item := range arr {
				if str, ok := item.(string); ok {
					newFavs = append(newFavs, FavouriteItem{UniqueId: str, Size: nil})
				} else {
					b, _ := json.Marshal(item)
					var fi FavouriteItem
					if err := json.Unmarshal(b, &fi); err == nil {
						newFavs = append(newFavs, fi)
					}
				}
			}
			user.Favourites = newFavs
		} else {
			b, _ := json.Marshal(favRaw)
			if err := json.Unmarshal(b, &user.Favourites); err != nil {
				user.Favourites = []FavouriteItem{}
			}
		}
	} else {
		user.Favourites = []FavouriteItem{}
	}

	if cartRaw, ok := raw["cart"]; ok && cartRaw != nil {
		b, _ := json.Marshal(cartRaw)
		if err := json.Unmarshal(b, &user.Cart); err != nil {
			user.Cart = []CartItem{}
		}
	} else {
		user.Cart = []CartItem{}
	}

	if ordersRaw, ok := raw["orders"]; ok && ordersRaw != nil {
		b, err := json.Marshal(ordersRaw)
		if err == nil {
			var orders []Order
			if err := json.Unmarshal(b, &orders); err == nil {
				user.Orders = orders
			} else {
				user.Orders = []Order{}
			}
		} else {
			user.Orders = []Order{}
		}
	} else {
		user.Orders = []Order{}
	}

	return user, nil
}

func updateUserInMokky(user *User, authHeader string) error {
	user.Password = ""
	body, _ := json.Marshal(user)
	url := fmt.Sprintf("https://3b7b2b24dfd8c527.mokky.dev/users/%d", user.ID)
	req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update failed with status %d", resp.StatusCode)
	}
	return nil
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	default:
		return ""
	}
}
