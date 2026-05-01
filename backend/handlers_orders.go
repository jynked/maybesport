package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

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

	rows, err := DB.Query(context.Background(), `
		SELECT oi.id, pi.unique_id, oi.size, oi.price, oi.quantity, oi.status,
		       pi.images, p.title_ru, p.title_en
		FROM order_items oi
		JOIN product_items pi ON oi.product_item_id = pi.id
		JOIN products p ON pi.product_id = p.id
		WHERE oi.order_id = $1
	`, orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type EnrichedOrderItemWithStatus struct {
		ID            int64                         `json:"id"`
		UniqueId      string                        `json:"uniqueId"`
		Size          interface{}                   `json:"size"`
		Price         int64                         `json:"price"`
		Quantity      int64                         `json:"quantity"`
		Status        OrderStatus                   `json:"status"`
		Title         Lang                          `json:"title"`
		Image         string                        `json:"image"`
		StatusHistory []OrderItemStatusHistoryEntry `json:"statusHistory,omitempty"`
	}
	enrichedItems := []EnrichedOrderItemWithStatus{}

	for rows.Next() {
		var it struct {
			ID       int64
			UniqueId string
			Size     interface{}
			Price    int64
			Quantity int64
			Status   OrderStatus
			Images   []string
			TitleRu  string
			TitleEn  string
		}
		err := rows.Scan(&it.ID, &it.UniqueId, &it.Size, &it.Price, &it.Quantity, &it.Status,
			&it.Images, &it.TitleRu, &it.TitleEn)
		if err != nil {
			continue
		}
		item := EnrichedOrderItemWithStatus{
			ID:       it.ID,
			UniqueId: it.UniqueId,
			Size:     it.Size,
			Price:    it.Price,
			Quantity: it.Quantity,
			Status:   it.Status,
			Title:    Lang{Ru: it.TitleRu, En: it.TitleEn},
			Image: func() string {
				if len(it.Images) > 0 {
					return it.Images[0]
				}
				return ""
			}(),
		}
		histRows, err := DB.Query(context.Background(), `
			SELECT status, description, timestamp FROM order_item_status_history WHERE order_item_id = $1 ORDER BY timestamp
		`, it.ID)
		if err == nil {
			var history []OrderItemStatusHistoryEntry
			for histRows.Next() {
				var h OrderItemStatusHistoryEntry
				histRows.Scan(&h.Status, &h.Description, &h.Timestamp)
				history = append(history, h)
			}
			histRows.Close()
			item.StatusHistory = history
		}
		enrichedItems = append(enrichedItems, item)
	}

	o.StatusHistory = []OrderStatusHistoryEntry{}

	response := struct {
		Order
		Items []EnrichedOrderItemWithStatus `json:"items"`
	}{
		Order: o,
		Items: enrichedItems,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
		CaptchaToken    string `json:"captchaToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if ok, err := verifyRecaptcha(req.CaptchaToken); !ok || err != nil {
		http.Error(w, "Captcha verification failed", http.StatusBadRequest)
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
		sizeStr := toString(it.Size)

		err = tx.QueryRow(context.Background(), `
        SELECT pi.id, pis.price, pis.quantity
        FROM product_items pi
        JOIN product_item_sizes pis ON pis.product_item_id = pi.id
        WHERE pi.unique_id = $1 AND pis.size = $2
        FOR UPDATE
    `, it.UniqueId, sizeStr).Scan(&productItemID, &price, &stock)
		if err != nil {
			tx.Rollback(context.Background())
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		if stock < it.Quantity {
			tx.Rollback(context.Background())
			http.Error(w, fmt.Sprintf("Not enough stock for %s size %v", it.UniqueId, it.Size), http.StatusConflict)
			return
		}

		res, err := tx.Exec(context.Background(), `
        UPDATE product_item_sizes
        SET quantity = quantity - $1
        WHERE product_item_id = $2 AND size = $3 AND quantity >= $1
    `, it.Quantity, productItemID, sizeStr)
		if err != nil {
			tx.Rollback(context.Background())
			http.Error(w, "Failed to update stock", http.StatusInternalServerError)
			return
		}
		affected := res.RowsAffected()
		if affected == 0 {
			tx.Rollback(context.Background())
			http.Error(w, "Stock changed during checkout", http.StatusConflict)
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
		var orderItemID int64
		err = tx.QueryRow(context.Background(), `
        INSERT INTO order_items (order_id, product_item_id, size, price, quantity)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, orderID, productItemID, toString(it.Size), it.Price, it.Quantity).Scan(&orderItemID)
		if err != nil {
			http.Error(w, "Failed to add order item", http.StatusInternalServerError)
			return
		}
		_, err = tx.Exec(context.Background(), `
        INSERT INTO order_item_status_history (order_item_id, status, description, timestamp)
        VALUES ($1, $2, $3, NOW())
    `, orderItemID, OrderStatusCreated, "Заказ оформлен")
		if err != nil {
			http.Error(w, "Failed to add order item status history", http.StatusInternalServerError)
			return
		}
	}

	for _, it := range orderItems {
		_, err = tx.Exec(context.Background(), `
        DELETE FROM user_cart 
        WHERE user_id = $1 
          AND product_item_id = (SELECT id FROM product_items WHERE unique_id = $2)
          AND size = $3
    `, userID, it.UniqueId, toString(it.Size))
		if err != nil {
			slog.Warn("Failed to delete cart item",
				"user_id", userID,
				"unique_id", it.UniqueId,
				"size", toString(it.Size),
				"error", err)
		}
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
