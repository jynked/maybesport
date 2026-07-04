package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const uploadDir = "uploads"

func init() {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}
}

func saveBase64Image(dataURL string) (string, error) {
	parts := strings.SplitN(dataURL, ",", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid data URL")
	}
	mimePart := parts[0]
	var ext string
	if strings.Contains(mimePart, "image/png") {
		ext = ".png"
	} else if strings.Contains(mimePart, "image/jpeg") || strings.Contains(mimePart, "image/jpg") {
		ext = ".jpg"
	} else {
		ext = ".png"
	}
	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	if len(decoded) > 5*1024*1024 {
		return "", fmt.Errorf("image too large (max 5MB)")
	}

	mimeType := http.DetectContentType(decoded)
	if mimeType != "image/png" && mimeType != "image/jpeg" && mimeType != "image/jpg" {
		return "", fmt.Errorf("invalid image format: only PNG and JPEG allowed")
	}

	filename := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, filename)
	if err := os.WriteFile(filePath, decoded, 0644); err != nil {
		return "", err
	}
	return "/uploads/" + filename, nil
}

func processImages(images []string, existingImages []string) ([]string, error) {
	result := make([]string, 0, len(images))
	for _, img := range images {
		if strings.HasPrefix(img, "data:image") {
			url, err := saveBase64Image(img)
			if err != nil {
				return nil, err
			}
			result = append(result, url)
		} else {
			result = append(result, img)
		}
	}
	oldMap := make(map[string]bool)
	for _, old := range existingImages {
		oldMap[old] = true
	}
	for _, new := range result {
		delete(oldMap, new)
	}
	for oldPath := range oldMap {
		safeRemoveImage(oldPath)
	}
	return result, nil
}

func (c *DBCache) AdminItemsHandler(w http.ResponseWriter, r *http.Request) {
	items := c.GetItems()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (c *DBCache) AdminItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	items := c.GetItems()
	for _, item := range items {
		if item.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.NotFound(w, r)
}

func (c *DBCache) CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	var input Item
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	tx, err := DB.Begin(context.Background())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(context.Background())

	var productID int
	err = tx.QueryRow(context.Background(), `
			INSERT INTO products (type_ru, type_en, brand, sport_id, category_ru, category_en, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`,
		input.Type.Ru, input.Type.En,
		input.Brand,
		input.SportID,
		input.Category.Ru, input.Category.En,
		time.Now(),
	).Scan(&productID)
	if err != nil {
		http.Error(w, "Failed to insert product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, st := range input.Structure {
		_, err = tx.Exec(context.Background(), `
			INSERT INTO product_structures (product_id, name_ru, name_en, percent)
			VALUES ($1, $2, $3, $4)
		`, productID, st.Name.Ru, st.Name.En, st.Percent)
		if err != nil {
			http.Error(w, "Failed to insert structure", http.StatusInternalServerError)
			return
		}
	}

	for idx, sub := range input.Items {
		processedImages, err := processImages(sub.Images, nil)
		if err != nil {
			http.Error(w, "Failed to process images: "+err.Error(), http.StatusInternalServerError)
			return
		}

		uniqueId := sub.UniqueId
		if uniqueId == "" {
			uniqueId = fmt.Sprintf("%d-%d", productID, idx+1)
		}
		var itemID int
		err = tx.QueryRow(context.Background(), `
			INSERT INTO product_items (product_id, unique_id, images, title_ru, title_en, description_ru, description_en)
			VALUES ($1, $2, $3, COALESCE($4, ''), COALESCE($5, ''), COALESCE($6, ''), COALESCE($7, ''))
			RETURNING id
		`, productID, uniqueId, processedImages, sub.Title.Ru, sub.Title.En, sub.Description.Ru, sub.Description.En).Scan(&itemID)
		if err != nil {
			slog.Error("Failed to insert product item", "error", err, "uniqueId", uniqueId, "productID", productID)
			http.Error(w, "Failed to insert product item: "+err.Error(), http.StatusInternalServerError)
			return
		}
		for _, col := range sub.Color {
			if col.Ru == "" && col.En == "" {
				continue
			}
			_, err = tx.Exec(context.Background(), `
				INSERT INTO product_item_colors (product_item_id, color_ru, color_en)
				VALUES ($1, $2, $3)
			`, itemID, col.Ru, col.En)
			if err != nil {
				slog.Error("Failed to insert color", "error", err, "color", col)
				http.Error(w, "Failed to insert color", http.StatusInternalServerError)
				return
			}
		}
		for _, tag := range sub.Tags {
			var tagID int
			err = tx.QueryRow(context.Background(), `
				INSERT INTO tags (name_ru, name_en) VALUES ($1, $2)
				ON CONFLICT (name_ru, name_en) DO UPDATE SET name_ru = EXCLUDED.name_ru
				RETURNING id
			`, tag.Ru, tag.En).Scan(&tagID)
			if err != nil {
				http.Error(w, "Failed to upsert tag", http.StatusInternalServerError)
				return
			}
			_, err = tx.Exec(context.Background(), `
				INSERT INTO product_item_tags (product_item_id, tag_id) VALUES ($1, $2)
				ON CONFLICT DO NOTHING
			`, itemID, tagID)
			if err != nil {
				http.Error(w, "Failed to link tag", http.StatusInternalServerError)
				return
			}
		}
		for _, sz := range sub.Sizes {
			_, err = tx.Exec(context.Background(), `
				INSERT INTO product_item_sizes (product_item_id, size, price_cny, is_on_request, quantity)
				VALUES ($1, $2, $3, $4, $5)
			`, itemID, toString(sz.Size), sz.PriceCny, sz.IsOnRequest, sz.Quantity)
			if err != nil {
				slog.Error("Failed to insert size", "error", err, "size", sz.Size, "product_item_id", itemID)
				http.Error(w, "Failed to insert size: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	adminID, _ := getUserIDFromToken(r)
	slog.Info("admin action",
		"action", "create_item",
		"admin_id", adminID,
		"product_id", productID)

	go c.Refresh()

	response := input
	response.ID = productID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (c *DBCache) UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var input Item
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	input.ID = productID

	tx, err := DB.Begin(context.Background())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `
			UPDATE products SET
				type_ru = $1, type_en = $2,
				brand = $3,
				sport_id = $4,
				category_ru = $5, category_en = $6,
				created_at = $7
			WHERE id = $8
		`,
		input.Type.Ru, input.Type.En,
		input.Brand,
		input.SportID,
		input.Category.Ru, input.Category.En,
		input.CreatedAt,
		productID,
	)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	oldImagesByItem := make(map[string][]string)
	rows, err := tx.Query(context.Background(), `SELECT unique_id, images FROM product_items WHERE product_id = $1`, productID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var uid string
			var imgs []string
			if err := rows.Scan(&uid, &imgs); err == nil {
				oldImagesByItem[uid] = imgs
			}
		}
	}

	_, err = tx.Exec(context.Background(), `
		DELETE FROM product_item_sizes WHERE product_item_id IN (SELECT id FROM product_items WHERE product_id = $1)
	`, productID)
	if err != nil {
		http.Error(w, "Failed to delete old sizes", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(context.Background(), `
		DELETE FROM product_item_tags WHERE product_item_id IN (SELECT id FROM product_items WHERE product_id = $1)
	`, productID)
	if err != nil {
		http.Error(w, "Failed to delete old tags", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(context.Background(), `
		DELETE FROM product_items WHERE product_id = $1
	`, productID)
	if err != nil {
		http.Error(w, "Failed to delete old product items", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM products WHERE id = $1`, productID)
	if err != nil {
		http.Error(w, "Failed to delete old product data", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(context.Background(), `
			INSERT INTO products (id, type_ru, type_en, brand, sport_id, category_ru, category_en, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, productID,
		input.Type.Ru, input.Type.En,
		input.Brand,
		input.SportID,
		input.Category.Ru, input.Category.En,
		input.CreatedAt,
	)
	if err != nil {
		http.Error(w, "Failed to re-insert product", http.StatusInternalServerError)
		return
	}

	for _, st := range input.Structure {
		_, err = tx.Exec(context.Background(), `
			INSERT INTO product_structures (product_id, name_ru, name_en, percent)
			VALUES ($1, $2, $3, $4)
		`, productID, st.Name.Ru, st.Name.En, st.Percent)
		if err != nil {
			http.Error(w, "Failed to insert structure", http.StatusInternalServerError)
			return
		}
	}

	for idx, sub := range input.Items {
		oldImgs := oldImagesByItem[sub.UniqueId]
		processedImages, err := processImages(sub.Images, oldImgs)
		if err != nil {
			http.Error(w, "Failed to process images", http.StatusInternalServerError)
			return
		}
		uniqueId := sub.UniqueId
		if uniqueId == "" {
			uniqueId = fmt.Sprintf("%d-%d", productID, idx+1)
		}
		var itemID int
		err = tx.QueryRow(context.Background(), `
			INSERT INTO product_items (product_id, unique_id, images, title_ru, title_en, description_ru, description_en)
			VALUES ($1, $2, $3, COALESCE($4, ''), COALESCE($5, ''), COALESCE($6, ''), COALESCE($7, ''))
			RETURNING id
		`, productID, uniqueId, processedImages, sub.Title.Ru, sub.Title.En, sub.Description.Ru, sub.Description.En).Scan(&itemID)
		if err != nil {
			slog.Error("Failed to insert product item", "error", err, "uniqueId", uniqueId, "productID", productID)
			http.Error(w, "Failed to insert product item: "+err.Error(), http.StatusInternalServerError)
			return
		}
		for _, tag := range sub.Tags {
			var tagID int
			err = tx.QueryRow(context.Background(), `
				INSERT INTO tags (name_ru, name_en) VALUES ($1, $2)
				ON CONFLICT (name_ru, name_en) DO UPDATE SET name_ru = EXCLUDED.name_ru
				RETURNING id
			`, tag.Ru, tag.En).Scan(&tagID)
			if err != nil {
				http.Error(w, "Failed to upsert tag", http.StatusInternalServerError)
				return
			}
			_, err = tx.Exec(context.Background(), `
				INSERT INTO product_item_tags (product_item_id, tag_id) VALUES ($1, $2)
				ON CONFLICT DO NOTHING
			`, itemID, tagID)
			if err != nil {
				http.Error(w, "Failed to link tag", http.StatusInternalServerError)
				return
			}
		}
		for _, col := range sub.Color {
			if col.Ru == "" && col.En == "" {
				continue
			}
			_, err = tx.Exec(context.Background(), `
				INSERT INTO product_item_colors (product_item_id, color_ru, color_en)
				VALUES ($1, $2, $3)
			`, itemID, col.Ru, col.En)
			if err != nil {
				slog.Error("Failed to insert color", "error", err, "color", col)
				http.Error(w, "Failed to insert color", http.StatusInternalServerError)
				return
			}
		}
		for _, sz := range sub.Sizes {
			_, err = tx.Exec(context.Background(), `
				INSERT INTO product_item_sizes (product_item_id, size, price_cny, is_on_request, quantity)
				VALUES ($1, $2, $3, $4, $5)
			`, itemID, toString(sz.Size), sz.PriceCny, sz.IsOnRequest, sz.Quantity)
			if err != nil {
				slog.Error("Failed to insert size", "error", err, "size", sz.Size, "product_item_id", itemID)
				http.Error(w, "Failed to insert size: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	adminID, _ := getUserIDFromToken(r)
	slog.Info("admin action",
		"action", "update_item",
		"admin_id", adminID,
		"product_id", productID)

	go c.Refresh()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(input)
}

func (c *DBCache) DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var images []string
	rows, err := DB.Query(context.Background(), `
        SELECT images FROM product_items WHERE product_id = $1
    `, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var imgs []string
			if err := rows.Scan(&imgs); err == nil {
				images = append(images, imgs...)
			}
		}
	}
	_, err = DB.Exec(context.Background(), "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, imgPath := range images {
		safeRemoveImage(imgPath)
	}
	adminID, _ := getUserIDFromToken(r)
	slog.Info("admin action",
		"action", "delete_item",
		"admin_id", adminID,
		"product_id", id)
	go c.Refresh()
	w.WriteHeader(http.StatusOK)
}

func (c *DBCache) AdminGetOrders(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	statusFilter := query.Get("status")
	userEmailFilter := query.Get("user_email")

	sql := `
        SELECT o.id, o.user_id, u.email, o.status, o.total_amount, o.delivery_address, o.created_at, o.updated_at
        FROM orders o
        JOIN users u ON o.user_id = u.id
        WHERE 1=1
    `
	countSql := `SELECT COUNT(*) FROM orders o JOIN users u ON o.user_id = u.id WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	if statusFilter != "" {
		sql += fmt.Sprintf(" AND o.status = $%d", argPos)
		countSql += fmt.Sprintf(" AND o.status = $%d", argPos)
		args = append(args, statusFilter)
		argPos++
	}
	if userEmailFilter != "" {
		sql += fmt.Sprintf(" AND u.email ILIKE $%d", argPos)
		countSql += fmt.Sprintf(" AND u.email ILIKE $%d", argPos)
		args = append(args, "%"+userEmailFilter+"%")
		argPos++
	}

	sql += fmt.Sprintf(" ORDER BY o.created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, limit, offset)

	var total int
	err := DB.QueryRow(context.Background(), countSql, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		http.Error(w, "Failed to count orders", http.StatusInternalServerError)
		return
	}

	rows, err := DB.Query(context.Background(), sql, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	orders := []AdminOrderListItem{}
	for rows.Next() {
		var o AdminOrderListItem
		err := rows.Scan(&o.ID, &o.UserID, &o.UserEmail, &o.Status, &o.TotalAmount, &o.DeliveryAddress, &o.CreatedAt, &o.UpdatedAt)
		if err != nil {
			continue
		}
		orders = append(orders, o)
	}

	response := map[string]interface{}{
		"orders":     orders,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + limit - 1) / limit,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *DBCache) AdminUpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Status      string `json:"status"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	validStatuses := map[string]bool{
		"created": true, "processing": true, "shipped": true,
		"delivered": true, "received": true, "cancelled": true,
	}
	if !validStatuses[req.Status] {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	tx, err := DB.Begin(context.Background())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `
        UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2
    `, req.Status, orderID)
	if err != nil {
		http.Error(w, "Failed to update order status", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(context.Background(), `
        INSERT INTO order_status_history (order_id, status, description, timestamp)
        VALUES ($1, $2, $3, NOW())
    `, orderID, req.Status, req.Description)
	if err != nil {
		http.Error(w, "Failed to add status history", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		http.Error(w, "Failed to commit", http.StatusInternalServerError)
		return
	}

	adminID, _ := getUserIDFromToken(r)
	slog.Info("admin action",
		"action", "update_order_status",
		"admin_id", adminID,
		"order_id", orderID,
		"new_status", req.Status)

	go c.Refresh()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Status updated"})
}

func (c *DBCache) AdminDeleteLastOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	tx, err := DB.Begin(context.Background())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(context.Background())

	var lastStatus string
	err = tx.QueryRow(context.Background(), `
        SELECT status FROM order_status_history 
        WHERE order_id = $1 
        ORDER BY timestamp DESC 
        LIMIT 1
    `, orderID).Scan(&lastStatus)
	if err != nil {
		http.Error(w, "No status history found", http.StatusNotFound)
		return
	}

	_, err = tx.Exec(context.Background(), `
        DELETE FROM order_status_history 
        WHERE id = (SELECT id FROM order_status_history WHERE order_id = $1 ORDER BY timestamp DESC LIMIT 1)
    `, orderID)
	if err != nil {
		http.Error(w, "Failed to delete status history", http.StatusInternalServerError)
		return
	}

	var previousStatus string
	err = tx.QueryRow(context.Background(), `
        SELECT status FROM order_status_history 
        WHERE order_id = $1 
        ORDER BY timestamp DESC 
        LIMIT 1
    `, orderID).Scan(&previousStatus)
	if err != nil {
		previousStatus = "created"
		_, err = tx.Exec(context.Background(), `
            INSERT INTO order_status_history (order_id, status, description, timestamp)
            VALUES ($1, $2, $3, NOW())
        `, orderID, previousStatus, "Начальный статус")
		if err != nil {
			http.Error(w, "Failed to insert default status", http.StatusInternalServerError)
			return
		}
	}

	_, err = tx.Exec(context.Background(), `
        UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2
    `, previousStatus, orderID)
	if err != nil {
		http.Error(w, "Failed to update order status", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		http.Error(w, "Failed to commit", http.StatusInternalServerError)
		return
	}

	go c.Refresh()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Last status removed"})
}

func (c *DBCache) AdminGetOrderDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var o Order
	var deliveryAddress *string
	err = DB.QueryRow(context.Background(), `
        SELECT id, created_at, updated_at, status, total_amount, delivery_address
        FROM orders WHERE id = $1
    `, orderID).Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt, &o.Status, &o.TotalAmount, &deliveryAddress)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	o.DeliveryAddress = deliveryAddress

	histRows, err := DB.Query(context.Background(), `
    SELECT status, description, timestamp FROM order_status_history WHERE order_id = $1 ORDER BY timestamp
`, orderID)
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
		o.StatusHistory = history
	}

	itemRows, err := DB.Query(context.Background(), `
        SELECT pi.unique_id, oi.size, oi.price, oi.quantity
        FROM order_items oi
        JOIN product_items pi ON oi.product_item_id = pi.id
        WHERE oi.order_id = $1
    `, orderID)
	if err == nil {
		var items []OrderItem
		for itemRows.Next() {
			var it OrderItem
			itemRows.Scan(&it.UniqueId, &it.Size, &it.Price, &it.Quantity)
			items = append(items, it)
		}
		itemRows.Close()
		o.Items = items
	}

	enriched := enrichOrderWithItems(&o, c.GetFlatItems())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enriched)
}

func (c *DBCache) AdminGetOrderItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	rows, err := DB.Query(context.Background(), `
        SELECT oi.id, pi.unique_id, oi.size, oi.price, oi.quantity, oi.status,
               pi.images, pi.title_ru, pi.title_en
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

	type ItemWithHistory struct {
		OrderItem
		Image         string                        `json:"image"`
		Title         Lang                          `json:"title"`
		StatusHistory []OrderItemStatusHistoryEntry `json:"statusHistory"`
	}
	var items []ItemWithHistory

	for rows.Next() {
		var it ItemWithHistory
		var images []string
		var titleRu, titleEn string
		err := rows.Scan(&it.OrderItem.ID, &it.OrderItem.UniqueId, &it.OrderItem.Size, &it.OrderItem.Price,
			&it.OrderItem.Quantity, &it.OrderItem.Status, &images, &titleRu, &titleEn)
		if err != nil {
			log.Printf("Scan error in AdminGetOrderItems: %v", err)
			continue
		}
		if len(images) > 0 {
			it.Image = images[0]
		}
		it.Title = Lang{Ru: titleRu, En: titleEn}

		histRows, err := DB.Query(context.Background(), `
            SELECT status, description, timestamp FROM order_item_status_history WHERE order_item_id = $1 ORDER BY timestamp
        `, it.OrderItem.ID)
		if err == nil {
			var history []OrderItemStatusHistoryEntry
			for histRows.Next() {
				var h OrderItemStatusHistoryEntry
				if err := histRows.Scan(&h.Status, &h.Description, &h.Timestamp); err == nil {
					history = append(history, h)
				}
			}
			histRows.Close()
			it.StatusHistory = history
		}
		items = append(items, it)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (c *DBCache) AdminUpdateOrderItemStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID, err := strconv.Atoi(vars["itemId"])
	if err != nil {
		http.Error(w, "Invalid order item ID", http.StatusBadRequest)
		return
	}
	var req struct {
		Status      string `json:"status"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	validStatuses := map[string]bool{
		"created": true, "processing": true, "shipped": true,
		"delivered": true, "received": true, "cancelled": true,
	}
	if !validStatuses[req.Status] {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	tx, err := DB.Begin(context.Background())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `
        UPDATE order_items SET status = $1 WHERE id = $2
    `, req.Status, itemID)
	if err != nil {
		http.Error(w, "Failed to update order item status", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(context.Background(), `
        INSERT INTO order_item_status_history (order_item_id, status, description, timestamp)
        VALUES ($1, $2, $3, NOW())
    `, itemID, req.Status, req.Description)
	if err != nil {
		http.Error(w, "Failed to add status history", http.StatusInternalServerError)
		return
	}

	var orderID int
	err = tx.QueryRow(context.Background(), "SELECT order_id FROM order_items WHERE id = $1", itemID).Scan(&orderID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	var allReceived bool
	var maxStatus string
	rows, _ := tx.Query(context.Background(), "SELECT status FROM order_items WHERE order_id = $1", orderID)
	statusPriority := map[string]int{
		"received": 5, "delivered": 4, "shipped": 3, "processing": 2, "created": 1, "cancelled": 0,
	}
	maxPriority := 0
	allReceived = true
	for rows.Next() {
		var s string
		rows.Scan(&s)
		if s != "received" {
			allReceived = false
		}
		if statusPriority[s] > maxPriority {
			maxPriority = statusPriority[s]
			maxStatus = s
		}
	}
	rows.Close()
	newOrderStatus := maxStatus
	if allReceived {
		newOrderStatus = "received"
	}
	_, err = tx.Exec(context.Background(), `
        UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2
    `, newOrderStatus, orderID)
	if err != nil {
		http.Error(w, "Failed to update order status", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		http.Error(w, "Failed to commit", http.StatusInternalServerError)
		return
	}

	adminID, _ := getUserIDFromToken(r)
	slog.Info("admin action",
		"action", "update_order_item_status",
		"admin_id", adminID,
		"order_item_id", itemID,
		"new_status", req.Status)

	go c.Refresh()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order item status updated"})
}

func safeRemoveImage(imagePath string) {
	if imagePath == "" {
		return
	}
	if !strings.HasPrefix(imagePath, "/uploads/") {
		return
	}
	filename := strings.TrimPrefix(imagePath, "/uploads/")
	if strings.Contains(filename, "..") || strings.ContainsAny(filename, "\\/:*?\"<>|") {
		log.Printf("SECURITY: Blocked suspicious file deletion attempt: %s", filename)
		return
	}
	fullPath := filepath.Join(uploadDir, filename)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to remove file %s: %v", fullPath, err)
	}
}

func (c *DBCache) AdminUploadHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(15 << 20)
	if err != nil {
		http.Error(w, "File too large or form parsing error", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		http.Error(w, "Cannot read file", http.StatusBadRequest)
		return
	}
	file.Seek(0, 0)

	mimeType := http.DetectContentType(buf)
	allowedImages := []string{"image/jpeg", "image/png", "image/gif"}
	allowedVideo := []string{"video/mp4", "video/webm", "video/quicktime"}

	isImage := false
	for _, t := range allowedImages {
		if mimeType == t {
			isImage = true
			break
		}
	}
	isVideo := false
	if !isImage {
		for _, t := range allowedVideo {
			if mimeType == t {
				isVideo = true
				break
			}
		}
	}
	if !isImage && !isVideo {
		http.Error(w, "Unsupported file type. Allowed: JPEG, PNG, GIF, MP4, WebM, MOV", http.StatusBadRequest)
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		if isImage {
			ext = ".png"
		} else {
			ext = ".gif"
		}
	}

	finalExt := ext
	if isVideo {
		finalExt = ".gif"
	}
	uniqueName := uuid.New().String() + finalExt
	tempPath := filepath.Join(uploadDir, "temp_"+uniqueName)
	finalPath := filepath.Join(uploadDir, uniqueName)

	outFile, err := os.Create(tempPath)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, file)
	if err != nil {
		os.Remove(tempPath)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	outFile.Close()

	if isVideo {

		cmd := exec.Command(
			"ffmpeg",
			"-i", tempPath,
			"-vf", "fps=60,scale=480:-1:flags=lanczos,split[s0][s1];[s0]palettegen[p];[s1][p]paletteuse",
			"-loop", "0",
			finalPath,
		)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			os.Remove(tempPath)
			slog.Error("ffmpeg conversion failed", "error", err, "stderr", stderr.String())
			http.Error(w, "Video conversion to GIF failed", http.StatusInternalServerError)
			return
		}

		os.Remove(tempPath)

		if _, err := os.Stat(finalPath); os.IsNotExist(err) {
			http.Error(w, "GIF not created", http.StatusInternalServerError)
			return
		}
	} else {

		if err := os.Rename(tempPath, finalPath); err != nil {
			os.Remove(tempPath)
			http.Error(w, "Failed to move file", http.StatusInternalServerError)
			return
		}
	}

	url := "/uploads/" + uniqueName
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": url})
}

func (c *DBCache) AdminUpdateMainPage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TitleRu         string `json:"title_ru"`
		TitleEn         string `json:"title_en"`
		Image           string `json:"image"`
		NewItemUniqueId string `json:"newItemUniqueId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Image == "" {
		http.Error(w, "image is required", http.StatusBadRequest)
		return
	}
	if req.NewItemUniqueId == "" {
		req.NewItemUniqueId = "1-1"
	}

	var currentImage string
	err := DB.QueryRow(context.Background(), "SELECT image FROM main_page_new WHERE id = 1").Scan(&currentImage)
	if err != nil && err != sql.ErrNoRows {
		slog.Error("failed to get current main page image", "error", err)
	}
	if currentImage != "" && strings.HasPrefix(currentImage, "/uploads/") && currentImage != req.Image {
		safeRemoveImage(currentImage)
		slog.Info("removed old banner image", "path", currentImage)
	}

	_, err = DB.Exec(context.Background(), `
        UPDATE main_page_new
        SET title_ru = $1, title_en = $2, image = $3, new_item_unique_id = $4
        WHERE id = 1
    `, req.TitleRu, req.TitleEn, req.Image, req.NewItemUniqueId)
	if err != nil {
		slog.Error("failed to update main page", "error", err)
		http.Error(w, "Failed to update main page", http.StatusInternalServerError)
		return
	}

	go c.Refresh()
	adminID, _ := getUserIDFromToken(r)
	slog.Info("admin updated main page", "admin_id", adminID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Main page updated successfully"})
}

func (c *DBCache) AdminGetSportsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(context.Background(), "SELECT id, name_ru, name_en FROM sports ORDER BY name_ru")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var sports []Sport
	for rows.Next() {
		var s Sport
		var nameRu, nameEn string
		if err := rows.Scan(&s.ID, &nameRu, &nameEn); err != nil {
			continue
		}
		s.Name = Lang{Ru: nameRu, En: nameEn}
		sports = append(sports, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sports)
}

func (c *DBCache) AdminGetUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	search := query.Get("search")

	sql := `SELECT id, email, name, created_at, deleted_at FROM users WHERE 1=1`
	countSql := `SELECT COUNT(*) FROM users WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	if search != "" {
		sql += fmt.Sprintf(" AND (email ILIKE $%d OR name ILIKE $%d)", argPos, argPos)
		countSql += fmt.Sprintf(" AND (email ILIKE $%d OR name ILIKE $%d)", argPos, argPos)
		args = append(args, "%"+search+"%")
		argPos++
	}

	sql += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, limit, offset)

	var total int
	err := DB.QueryRow(context.Background(), countSql, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		http.Error(w, "Failed to count users", http.StatusInternalServerError)
		return
	}

	rows, err := DB.Query(context.Background(), sql, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type UserListItem struct {
		ID        int        `json:"id"`
		Email     string     `json:"email"`
		Name      string     `json:"name"`
		CreatedAt time.Time  `json:"created_at"`
		DeletedAt *time.Time `json:"deleted_at"`
		IsBanned  bool       `json:"is_banned"`
	}
	users := []UserListItem{}
	for rows.Next() {
		var u UserListItem
		var deletedAt *time.Time
		err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &deletedAt)
		if err != nil {
			continue
		}
		if deletedAt != nil {
			u.DeletedAt = deletedAt
			u.IsBanned = true
		}
		users = append(users, u)
	}

	response := map[string]interface{}{
		"users":      users,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + limit - 1) / limit,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *DBCache) AdminBanUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Ban bool `json:"ban"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var query string
	if req.Ban {
		query = "UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL"
	} else {
		query = "UPDATE users SET deleted_at = NULL WHERE id = $1"
	}
	_, err = DB.Exec(context.Background(), query, userID)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User status updated"})
}
