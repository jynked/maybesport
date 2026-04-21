package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (c *DBCache) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	rows, err := DB.Query(context.Background(), `
		SELECT id, created_at, updated_at, status, total_amount, delivery_address
		FROM orders WHERE user_id = $1 ORDER BY created_at DESC
	`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var orders []Order
	for rows.Next() {
		var o Order
		var deliveryAddress *string
		err := rows.Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt, &o.Status, &o.TotalAmount, &deliveryAddress)
		if err != nil {
			continue
		}
		o.DeliveryAddress = deliveryAddress
		orders = append(orders, o)
	}

	for i, ord := range orders {
		histRows, err := DB.Query(context.Background(), `
			SELECT status, description, timestamp FROM order_status_history WHERE order_id = $1 ORDER BY timestamp
		`, ord.ID)
		if err == nil {
			var history []OrderStatusHistoryEntry
			for histRows.Next() {
				var h OrderStatusHistoryEntry
				err := histRows.Scan(&h.Status, &h.Description, &h.Timestamp)
				if err != nil {
					continue
				}
				history = append(history, h)
			}
			histRows.Close()
			orders[i].StatusHistory = history
		}

		itemRows, _ := DB.Query(context.Background(), `
			SELECT pi.unique_id, oi.size, oi.price, oi.quantity
			FROM order_items oi
			JOIN product_items pi ON oi.product_item_id = pi.id
			WHERE oi.order_id = $1
		`, ord.ID)
		var items []OrderItem
		for itemRows.Next() {
			var it OrderItem
			itemRows.Scan(&it.UniqueId, &it.Size, &it.Price, &it.Quantity)
			items = append(items, it)
		}
		itemRows.Close()
		orders[i].Items = items
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (c *DBCache) GetOrderDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var o Order
	var deliveryAddress *string
	err = DB.QueryRow(context.Background(), `
		SELECT id, created_at, updated_at, status, total_amount, delivery_address
		FROM orders WHERE id=$1 AND user_id=$2
	`, orderID, userID).Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt, &o.Status, &o.TotalAmount, &deliveryAddress)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	o.DeliveryAddress = deliveryAddress

	histRows, _ := DB.Query(context.Background(), `
		SELECT status, description, timestamp FROM order_status_history WHERE order_id = $1 ORDER BY timestamp
	`, orderID)
	var history []OrderStatusHistoryEntry
	for histRows.Next() {
		var h OrderStatusHistoryEntry
		var ts string
		histRows.Scan(&h.Status, &h.Description, &h.Timestamp)
		h.Timestamp, _ = time.Parse(time.RFC3339, ts)
		history = append(history, h)
	}
	histRows.Close()
	o.StatusHistory = history

	itemRows, _ := DB.Query(context.Background(), `
		SELECT pi.unique_id, oi.size, oi.price, oi.quantity
		FROM order_items oi
		JOIN product_items pi ON oi.product_item_id = pi.id
		WHERE oi.order_id = $1
	`, orderID)
	var items []OrderItem
	for itemRows.Next() {
		var it OrderItem
		itemRows.Scan(&it.UniqueId, &it.Size, &it.Price, &it.Quantity)
		items = append(items, it)
	}
	itemRows.Close()
	o.Items = items

	enriched := enrichOrderWithItems(&o, c.GetFlatItems())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enriched)
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
	enriched := EnrichedOrder{Order: *order}
	itemMap := make(map[string]ItemFlatten)
	for _, it := range flatItems {
		itemMap[it.UniqueId] = it
	}
	for _, oi := range order.Items {
		eoi := EnrichedOrderItem{OrderItem: oi}
		if flat, ok := itemMap[oi.UniqueId]; ok {
			eoi.Title = flat.Title
			if len(flat.Images) > 0 {
				eoi.Image = flat.Images[0]
			}
		}
		enriched.Items = append(enriched.Items, eoi)
	}
	return enriched
}

func (c *DBCache) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Items []struct {
			UniqueId string      `json:"uniqueId"`
			Size     interface{} `json:"size"`
			Quantity int64       `json:"quantity"`
		} `json:"items"`
		DeliveryAddress string `json:"deliveryAddress"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if len(req.Items) == 0 {
		http.Error(w, "No items in order", http.StatusBadRequest)
		return
	}

	tx, err := DB.Begin(context.Background())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(context.Background())

	var totalAmount int64
	orderItems := []OrderItem{}
	for _, it := range req.Items {
		var price int64
		var productItemID int
		var stock int64
		err = tx.QueryRow(context.Background(), `
            SELECT pi.id, pis.price, pis.quantity
            FROM product_items pi
            JOIN product_item_sizes pis ON pis.product_item_id = pi.id
            WHERE pi.unique_id = $1 AND pis.size = $2
        `, it.UniqueId, toString(it.Size)).Scan(&productItemID, &price, &stock)
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		if stock < it.Quantity {
			http.Error(w, fmt.Sprintf("Not enough stock for item %s size %v", it.UniqueId, it.Size), http.StatusConflict)
			return
		}
		totalAmount += price * it.Quantity
		orderItems = append(orderItems, OrderItem{
			UniqueId: it.UniqueId,
			Size:     it.Size,
			Price:    price,
			Quantity: it.Quantity,
		})
	}

	var orderID int
	err = tx.QueryRow(context.Background(), `
        INSERT INTO orders (user_id, status, total_amount, delivery_address, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id
    `, userID, OrderStatusCreated, totalAmount, req.DeliveryAddress).Scan(&orderID)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(context.Background(), `
        INSERT INTO order_status_history (order_id, status, description, timestamp)
        VALUES ($1, $2, $3, NOW())
    `, orderID, OrderStatusCreated, "Заказ оформлен")
	if err != nil {
		http.Error(w, "Failed to create status history", http.StatusInternalServerError)
		return
	}

	for _, it := range orderItems {
		var productItemID int
		err = tx.QueryRow(context.Background(), `
            SELECT id FROM product_items WHERE unique_id = $1
        `, it.UniqueId).Scan(&productItemID)
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		_, err = tx.Exec(context.Background(), `
            INSERT INTO order_items (order_id, product_item_id, size, price, quantity)
            VALUES ($1, $2, $3, $4, $5)
        `, orderID, productItemID, toString(it.Size), it.Price, it.Quantity)
		if err != nil {
			http.Error(w, "Failed to add order item", http.StatusInternalServerError)
			return
		}
		_, err = tx.Exec(context.Background(), `
            UPDATE product_item_sizes
            SET quantity = quantity - $1
            WHERE product_item_id = $2 AND size = $3
        `, it.Quantity, productItemID, toString(it.Size))
		if err != nil {
			http.Error(w, "Failed to update stock", http.StatusInternalServerError)
			return
		}
	}

	_, err = tx.Exec(context.Background(), `
        DELETE FROM user_cart WHERE user_id = $1
    `, userID)
	if err != nil {
		http.Error(w, "Failed to clear cart", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		http.Error(w, "Failed to commit order", http.StatusInternalServerError)
		return
	}

	go c.Refresh()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"orderId": orderID,
		"message": "Order created successfully",
	})
}
