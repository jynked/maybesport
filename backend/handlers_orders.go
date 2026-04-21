package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (c *Cache) GetUserOrders(w http.ResponseWriter, r *http.Request) {
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
	if fullUser.Orders == nil {
		fullUser.Orders = []Order{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fullUser.Orders)
}

func (c *Cache) GetOrderDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDStr := vars["id"]
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
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

	var targetOrder *Order
	for i, order := range fullUser.Orders {
		if order.ID == orderID {
			targetOrder = &fullUser.Orders[i]
			break
		}
	}
	if targetOrder == nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	enrichedOrder := enrichOrderWithItems(targetOrder, c.GetFlatItems())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrichedOrder)
}

type EnrichedOrderItem struct {
	OrderItem
	Title Lang   `json:"title"`
	Image string `json:"image"`
}

type EnrichedOrder struct {
	Order
	Items []EnrichedOrderItem `json:"items"`
}

func enrichOrderWithItems(order *Order, flatItems []ItemFlatten) EnrichedOrder {
	enriched := EnrichedOrder{
		Order: *order,
		Items: make([]EnrichedOrderItem, 0, len(order.Items)),
	}
	itemMap := make(map[string]ItemFlatten)
	for _, it := range flatItems {
		itemMap[it.UniqueId] = it
	}

	for _, orderItem := range order.Items {
		ei := EnrichedOrderItem{OrderItem: orderItem}
		if flat, ok := itemMap[orderItem.UniqueId]; ok {
			ei.Title = flat.Title
			if len(flat.Images) > 0 {
				ei.Image = flat.Images[0]
			}
		}
		enriched.Items = append(enriched.Items, ei)
	}
	return enriched
}
