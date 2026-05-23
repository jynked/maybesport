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
		SELECT id, type_ru, type_en, brand, COALESCE(sport_id, 0) as sport_id, category_ru, category_en, created_at
		FROM products
	`)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	productsMap := make(map[int]*Item)
	for rows.Next() {
		var p Item
		var typeRu, typeEn, brand, categoryRu, categoryEn string
		var createdAt time.Time
		var sportID int
		err := rows.Scan(&p.ID, &typeRu, &typeEn, &brand, &sportID, &categoryRu, &categoryEn, &createdAt)
		if err != nil {
			return nil, nil, err
		}
		p.Type = Lang{Ru: typeRu, En: typeEn}
		p.Brand = brand
		p.SportID = sportID
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
		SELECT pi.id, pi.product_id, pi.unique_id, pi.images,
		COALESCE(pi.title_ru, '') as title_ru,
		COALESCE(pi.title_en, '') as title_en,
		COALESCE(pi.description_ru, '') as description_ru,
		COALESCE(pi.description_en, '') as description_en,
		COALESCE((SELECT json_agg(json_build_object('size', size, 'isOnRequest', is_on_request, 'quantity', quantity, 'priceCny', price_cny)) 
			FROM product_item_sizes WHERE product_item_id = pi.id), '[]') as sizes_json,
		COALESCE((SELECT json_agg(json_build_object('ru', t.name_ru, 'en', t.name_en)) 
			FROM product_item_tags pit JOIN tags t ON pit.tag_id = t.id WHERE pit.product_item_id = pi.id), '[]') as tags_json,
		COALESCE((SELECT json_agg(json_build_object('ru', c.color_ru, 'en', c.color_en)) 
			FROM product_item_colors c WHERE c.product_item_id = pi.id), '[]') as colors_json
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
		var titleRu, titleEn, descRu, descEn string
		var sizesJSON, tagsJSON, colorsJSON string
		err := itemRows.Scan(&id, &productID, &uniqueID, &images,
			&titleRu, &titleEn, &descRu, &descEn,
			&sizesJSON, &tagsJSON, &colorsJSON)
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
		var colors []Lang
		json.Unmarshal([]byte(colorsJSON), &colors)

		sub := SubItem{
			UniqueId:    uniqueID,
			Title:       Lang{Ru: titleRu, En: titleEn},
			Description: Lang{Ru: descRu, En: descEn},
			Images:      images,
			Color:       colors,
			Tags:        tags,
			Sizes:       sizes,
		}
		parent.Items = append(parent.Items, sub)

		// Получаем название спорта (отдельным запросом или JOIN – сделаем отдельно)
		sportName := Lang{}
		if parent.SportID > 0 {
			var sportRu, sportEn string
			err = DB.QueryRow(context.Background(),
				"SELECT name_ru, name_en FROM sports WHERE id = $1", parent.SportID).Scan(&sportRu, &sportEn)
			if err == nil {
				sportName = Lang{Ru: sportRu, En: sportEn}
			}
		}

		flatten := ItemFlatten{
			UniqueId:    uniqueID,
			ID:          parent.ID,
			Type:        parent.Type,
			Title:       Lang{Ru: titleRu, En: titleEn},
			Description: Lang{Ru: descRu, En: descEn},
			Brand:       parent.Brand,
			Structure:   parent.Structure,
			Category:    parent.Category,
			CreatedAt:   parent.CreatedAt,
			Images:      images,
			Color:       colors,
			Tags:        tags,
			Sizes:       sizes,
			Sport:       sportName,
		}

		var totalQty int64
		var minPriceCny int64 = 1<<62 - 1
		var hasAvailable bool
		var hasOnRequest bool
		for _, sz := range sizes {
			if sz.PriceCny < minPriceCny {
				minPriceCny = sz.PriceCny
			}
			totalQty += sz.Quantity
			if sz.Quantity > 0 && !sz.IsOnRequest {
				hasAvailable = true
			}
			if sz.IsOnRequest {
				hasOnRequest = true
			}
		}
		if minPriceCny == 1<<62-1 {
			minPriceCny = 0
		}
		rate := GetCurrentExchangeRate()
		flatten.MinPrice = int64(float64(minPriceCny) * rate)
		flatten.TotalQuantity = totalQty
		if hasAvailable {
			flatten.Availability = "available"
		} else if hasOnRequest {
			flatten.Availability = "on_request"
		} else {
			flatten.Availability = "out_of_stock"
		}
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
	err := DB.QueryRow(context.Background(), "SELECT id, title_ru, title_en, image, new_item_unique_id FROM main_page_new LIMIT 1").
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
