package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"
)

type DBCache struct {
	mu          sync.RWMutex
	flatItems   []ItemFlatten
	items       []Item
	mainPageNew interface{}
}

func NewDBCache() *DBCache {
	return &DBCache{}
}

func (c *DBCache) Refresh() error {
	flatItems, items, err := loadItemsFromDB()
	if err != nil {
		return err
	}
	mainPage, err := loadMainPageFromDB()
	if err != nil {
		log.Printf("Warning: failed to load main page: %v", err)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.flatItems = flatItems
	c.items = items
	c.mainPageNew = mainPage
	return nil
}

func (c *DBCache) GetFlatItems() []ItemFlatten {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.flatItems
}

func (c *DBCache) GetItems() []Item {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.items
}

func (c *DBCache) GetMainPageNew() interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.mainPageNew
}

func loadItemsFromDB() ([]ItemFlatten, []Item, error) {
	rows, err := DB.Query(context.Background(), `
        SELECT id, type_ru, type_en, title_ru, title_en, desc_ru, desc_en,
               brand, country_ru, country_en, category_ru, category_en, created_at
        FROM products
    `)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	productsMap := make(map[int]*Item)
	for rows.Next() {
		var p Item
		var typeRu, typeEn, titleRu, titleEn, descRu, descEn, brand, countryRu, countryEn, categoryRu, categoryEn string
		var createdAt time.Time
		err := rows.Scan(&p.ID, &typeRu, &typeEn, &titleRu, &titleEn, &descRu, &descEn,
			&brand, &countryRu, &countryEn, &categoryRu, &categoryEn, &createdAt)
		if err != nil {
			return nil, nil, err
		}
		p.Type = Lang{Ru: typeRu, En: typeEn}
		p.Title = Lang{Ru: titleRu, En: titleEn}
		p.Description = Lang{Ru: descRu, En: descEn}
		p.Brand = brand
		p.Country = Lang{Ru: countryRu, En: countryEn}
		p.Category = Lang{Ru: categoryRu, En: categoryEn}
		p.CreatedAt = createdAt
		p.Items = []SubItem{}
		productsMap[p.ID] = &p
	}

	structRows, err := DB.Query(context.Background(), "SELECT product_id, name_ru, name_en, percent FROM product_structures")
	if err != nil {
		return nil, nil, err
	}
	defer structRows.Close()
	for structRows.Next() {
		var productID int
		var nameRu, nameEn string
		var percent int
		if err := structRows.Scan(&productID, &nameRu, &nameEn, &percent); err != nil {
			continue
		}
		if p, ok := productsMap[productID]; ok {
			p.Structure = append(p.Structure, struct {
				Name    Lang `json:"name"`
				Percent int  `json:"percent"`
			}{
				Name:    Lang{Ru: nameRu, En: nameEn},
				Percent: percent,
			})
		}
	}

	itemRows, err := DB.Query(context.Background(), `
		SELECT pi.id, pi.product_id, pi.unique_id, pi.images, pi.color_ru, pi.color_en,
		       COALESCE((SELECT json_agg(json_build_object('size', size, 'price', price, 'isOnRequest', is_on_request, 'quantity', quantity, 'priceCny', price_cny)) 
		                 FROM product_item_sizes WHERE product_item_id = pi.id), '[]') as sizes_json,
		       COALESCE((SELECT json_agg(json_build_object('ru', t.name_ru, 'en', t.name_en)) 
		                 FROM product_item_tags pit JOIN tags t ON pit.tag_id = t.id WHERE pit.product_item_id = pi.id), '[]') as tags_json
		FROM product_items pi
	`)
	if err != nil {
		return nil, nil, err
	}
	defer itemRows.Close()

	var flatItems []ItemFlatten
	for itemRows.Next() {
		var id, productID int
		var uniqueID string
		var images []string
		var colorRu, colorEn string
		var sizesJSON, tagsJSON string
		err := itemRows.Scan(&id, &productID, &uniqueID, &images, &colorRu, &colorEn, &sizesJSON, &tagsJSON)
		if err != nil {
			return nil, nil, err
		}
		parent := productsMap[productID]
		if parent == nil {
			continue
		}
		var sizes []Size
		json.Unmarshal([]byte(sizesJSON), &sizes)
		var tags []Lang
		json.Unmarshal([]byte(tagsJSON), &tags)

		sub := SubItem{
			UniqueId: uniqueID,
			Images:   images,
			Color:    []Lang{{Ru: colorRu, En: colorEn}},
			Tags:     tags,
			Sizes:    sizes,
		}
		parent.Items = append(parent.Items, sub)

		flatten := ItemFlatten{
			UniqueId:    uniqueID,
			ID:          parent.ID,
			Type:        parent.Type,
			Title:       parent.Title,
			Description: parent.Description,
			Brand:       parent.Brand,
			Country:     parent.Country,
			Structure:   parent.Structure,
			Category:    parent.Category,
			CreatedAt:   parent.CreatedAt,
			Images:      images,
			Color:       []Lang{{Ru: colorRu, En: colorEn}},
			Tags:        tags,
			Sizes:       sizes,
		}
		var minPrice int64 = 1 << 62
		var totalQty int64
		avail := "out_of_stock"
		for _, sz := range sizes {
			if sz.Price < minPrice {
				minPrice = sz.Price
			}
			totalQty += sz.Quantity
			if sz.Quantity > 0 && !sz.IsOnRequest {
				avail = "available"
			} else if sz.Quantity > 0 && sz.IsOnRequest && avail != "available" {
				avail = "on_request"
			}
		}
		flatten.MinPrice = minPrice
		flatten.TotalQuantity = totalQty
		flatten.Availability = avail
		flatItems = append(flatItems, flatten)
	}

	for i, flat := range flatItems {
		var siblings []SiblingItem
		for _, other := range flatItems {
			if other.ID == flat.ID && other.UniqueId != flat.UniqueId {
				img := ""
				if len(other.Images) > 0 {
					img = other.Images[0]
				}
				siblings = append(siblings, SiblingItem{UniqueId: other.UniqueId, Image: img})
			}
		}
		flatItems[i].SiblingItems = siblings
	}

	itemsSlice := make([]Item, 0, len(productsMap))
	for _, v := range productsMap {
		itemsSlice = append(itemsSlice, *v)
	}
	return flatItems, itemsSlice, nil
}

func loadMainPageFromDB() (interface{}, error) {
	var titleRu, titleEn, image string
	var productItemUniqueID *string
	var id int
	err := DB.QueryRow(context.Background(), "SELECT id, title_ru, title_en, image, product_item_unique_id FROM main_page_new LIMIT 1").
		Scan(&id, &titleRu, &titleEn, &image, &productItemUniqueID)
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{
		"id":    id,
		"title": map[string]string{"ru": titleRu, "en": titleEn},
		"image": image,
	}
	if productItemUniqueID != nil {
		result["newItemUniqueId"] = *productItemUniqueID
	}
	return result, nil
}
