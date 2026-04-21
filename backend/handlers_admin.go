package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
		if strings.HasPrefix(oldPath, "/uploads/") {
			filename := strings.TrimPrefix(oldPath, "/uploads/")
			fullPath := filepath.Join(uploadDir, filename)
			os.Remove(fullPath)
		}
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
		INSERT INTO products (type_ru, type_en, title_ru, title_en, desc_ru, desc_en,
			brand, country_ru, country_en, category_ru, category_en, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`,
		input.Type.Ru, input.Type.En,
		input.Title.Ru, input.Title.En,
		input.Description.Ru, input.Description.En,
		input.Brand,
		input.Country.Ru, input.Country.En,
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
		colorRu, colorEn := "", ""
		if len(sub.Color) > 0 {
			colorRu = sub.Color[0].Ru
			colorEn = sub.Color[0].En
		}
		uniqueId := sub.UniqueId
		if uniqueId == "" {
			uniqueId = fmt.Sprintf("%d-%d", productID, idx+1)
		}
		var itemID int
		err = tx.QueryRow(context.Background(), `
			INSERT INTO product_items (product_id, unique_id, images, color_ru, color_en)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, productID, uniqueId, processedImages, colorRu, colorEn).Scan(&itemID)
		if err != nil {
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
		for _, sz := range sub.Sizes {
			_, err = tx.Exec(context.Background(), `
				INSERT INTO product_item_sizes (product_item_id, size, price, is_on_request, quantity, price_cny)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, itemID, toString(sz.Size), sz.Price, sz.IsOnRequest, sz.Quantity, sz.PriceCny)
			if err != nil {
				http.Error(w, "Failed to insert size", http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

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
			title_ru = $3, title_en = $4,
			desc_ru = $5, desc_en = $6,
			brand = $7,
			country_ru = $8, country_en = $9,
			category_ru = $10, category_en = $11,
			created_at = $12
		WHERE id = $13
	`,
		input.Type.Ru, input.Type.En,
		input.Title.Ru, input.Title.En,
		input.Description.Ru, input.Description.En,
		input.Brand,
		input.Country.Ru, input.Country.En,
		input.Category.Ru, input.Category.En,
		input.CreatedAt, productID,
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

	_, err = tx.Exec(context.Background(), `DELETE FROM products WHERE id = $1`, productID)
	if err != nil {
		http.Error(w, "Failed to delete old product data", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(context.Background(), `
		INSERT INTO products (id, type_ru, type_en, title_ru, title_en, desc_ru, desc_en,
			brand, country_ru, country_en, category_ru, category_en, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, productID,
		input.Type.Ru, input.Type.En,
		input.Title.Ru, input.Title.En,
		input.Description.Ru, input.Description.En,
		input.Brand,
		input.Country.Ru, input.Country.En,
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
		colorRu, colorEn := "", ""
		if len(sub.Color) > 0 {
			colorRu = sub.Color[0].Ru
			colorEn = sub.Color[0].En
		}
		uniqueId := sub.UniqueId
		if uniqueId == "" {
			uniqueId = fmt.Sprintf("%d-%d", productID, idx+1)
		}
		var itemID int
		err = tx.QueryRow(context.Background(), `
			INSERT INTO product_items (product_id, unique_id, images, color_ru, color_en)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, productID, uniqueId, processedImages, colorRu, colorEn).Scan(&itemID)
		if err != nil {
			http.Error(w, "Failed to insert product item", http.StatusInternalServerError)
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
		for _, sz := range sub.Sizes {
			_, err = tx.Exec(context.Background(), `
				INSERT INTO product_item_sizes (product_item_id, size, price, is_on_request, quantity, price_cny)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, itemID, toString(sz.Size), sz.Price, sz.IsOnRequest, sz.Quantity, sz.PriceCny)
			if err != nil {
				http.Error(w, "Failed to insert size", http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

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
		if strings.HasPrefix(imgPath, "/uploads/") {
			filename := strings.TrimPrefix(imgPath, "/uploads/")
			fullPath := filepath.Join(uploadDir, filename)
			os.Remove(fullPath)
		}
	}
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
