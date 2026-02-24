package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Cache struct {
	mu          sync.RWMutex
	items       []Item
	flatItems   []ItemFlatten
	mainPageNew interface{}
	lastUpdated time.Time
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) Refresh() error {
	if err := c.loadItems(); err != nil {
		return err
	}
	if err := c.loadMainPageNew(); err != nil {
		log.Printf("Warning: failed to load main_page_new: %v", err)
	}
	return nil
}

func (c *Cache) loadItems() error {
	resp, err := http.Get("https://3b7b2b24dfd8c527.mokky.dev/items")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var items []Item
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return err
	}

	flat := flattenItems(items)

	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = items
	c.flatItems = flat
	c.lastUpdated = time.Now()
	return nil
}

func (c *Cache) loadMainPageNew() error {
	resp, err := http.Get("https://3b7b2b24dfd8c527.mokky.dev/main_page_new")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var arr []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&arr); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if len(arr) > 0 {
		c.mainPageNew = arr[0]
	} else {
		c.mainPageNew = nil
	}
	return nil
}

func flattenItems(items []Item) []ItemFlatten {
	var flat []ItemFlatten
	siblingMap := make(map[int][]SiblingItem)

	for _, p := range items {
		var siblings []SiblingItem
		for _, sub := range p.Items {
			image := ""
			if len(sub.Images) > 0 {
				image = sub.Images[0]
			}
			sibling := SiblingItem{
				UniqueId: sub.UniqueId,
				Image:    image,
			}
			siblings = append(siblings, sibling)
		}
		siblingMap[p.ID] = siblings
	}

	for _, p := range items {
		allSiblings := siblingMap[p.ID]
		for _, sub := range p.Items {
			minPrice := int(^uint(0) >> 1)
			totalQty := 0
			avail := "out_of_stock"

			for _, s := range sub.Sizes {
				if s.Price < minPrice {
					minPrice = s.Price
				}
				totalQty += s.Quantity
				if s.Quantity > 0 && !s.IsOnRequest {
					avail = "available"
				} else if s.Quantity > 0 && s.IsOnRequest && avail != "available" {
					avail = "on_request"
				}
			}

			f := ItemFlatten{
				UniqueId:      sub.UniqueId,
				ID:            p.ID,
				Type:          p.Type,
				Title:         p.Title,
				Description:   p.Description,
				Brand:         p.Brand,
				Country:       p.Country,
				Structure:     p.Structure,
				Category:      p.Category,
				CreatedAt:     p.CreatedAt,
				Images:        sub.Images,
				Color:         sub.Color,
				Tags:          sub.Tags,
				Sizes:         sub.Sizes,
				MinPrice:      minPrice,
				TotalQuantity: totalQty,
				Availability:  avail,
			}

			var others []SiblingItem
			for _, sib := range allSiblings {
				if sib.UniqueId != sub.UniqueId {
					others = append(others, sib)
				}
			}
			f.SiblingItems = others

			flat = append(flat, f)
		}
	}
	return flat
}

func (c *Cache) GetItems() []Item {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.items
}

func (c *Cache) GetFlatItems() []ItemFlatten {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.flatItems
}

func (c *Cache) GetMainPageNew() interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.mainPageNew
}
