package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
)

func (c *Cache) ItemsHandler(w http.ResponseWriter, r *http.Request) {
	params := parseFilterParams(r)
	items := c.GetFlatItems()
	filtered := filterItems(items, params)

	start := (params.Page - 1) * params.Limit
	if start < 0 {
		start = 0
	}
	end := start + params.Limit
	if end > len(filtered) {
		end = len(filtered)
	}
	if start > len(filtered) {
		start = len(filtered)
		end = len(filtered)
	}
	paginated := filtered[start:end]

	applied := buildAppliedFilters(params)

	response := ItemsResponse{
		Items:          paginated,
		Total:          len(filtered),
		Page:           params.Page,
		Limit:          params.Limit,
		TotalPages:     (len(filtered) + params.Limit - 1) / params.Limit,
		AppliedFilters: applied,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func buildAppliedFilters(params FilterParams) []AppliedFilter {
	var filters []AppliedFilter

	if params.PriceMin != nil || params.PriceMax != nil {
		minVal := 0
		if params.PriceMin != nil {
			minVal = *params.PriceMin
		}
		maxVal := "∞"
		if params.PriceMax != nil {
			maxVal = strconv.Itoa(*params.PriceMax)
		}
		filters = append(filters, AppliedFilter{
			Key:   "price",
			Value: fmt.Sprintf("%d - %s", minVal, maxVal),
			Label: "Цена",
		})
	}

	addFilter := func(key, label string, values []string) {
		if len(values) > 0 {
			filters = append(filters, AppliedFilter{
				Key:   key,
				Value: values,
				Label: label,
			})
		}
	}

	addFilter("brands", "Бренд", params.Brands)
	addFilter("countries", "Страна", params.Countries)
	addFilter("materials", "Материал", params.Materials)
	addFilter("categories", "Категория", params.Categories)
	addFilter("types", "Тип", params.Types)
	addFilter("colors", "Цвет", params.Colors)
	addFilter("tags", "Тег", params.Tags)
	addFilter("sizes", "Размер", params.Sizes)
	addFilter("availability", "Наличие", params.Availability)

	return filters
}

func (c *Cache) ItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uniqueId := vars["uniqueId"]

	items := c.GetFlatItems()
	for _, item := range items {
		if item.UniqueId == uniqueId {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.NotFound(w, r)
}

func (c *Cache) SimilarHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uniqueId := vars["uniqueId"]

	limitStr := r.URL.Query().Get("limit")
	limit := 4
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	items := c.GetFlatItems()

	var current *ItemFlatten
	for _, item := range items {
		if item.UniqueId == uniqueId {
			current = &item
			break
		}
	}
	if current == nil {
		http.NotFound(w, r)
		return
	}

	var candidates []ItemFlatten
	for _, item := range items {
		if item.ID == current.ID {
			continue
		}
		candidates = append(candidates, item)
	}

	type scoredItem struct {
		Item  ItemFlatten
		Score int
	}
	scored := make([]scoredItem, len(candidates))
	for i, cand := range candidates {
		score := 0

		if cand.Category.Ru == current.Category.Ru || cand.Category.En == current.Category.En {
			score += 30
		}

		if cand.Brand == current.Brand {
			score += 25
		}

		if cand.Type.Ru == current.Type.Ru || cand.Type.En == current.Type.En {
			score += 20
		}

		priceDiff := abs(cand.MinPrice - current.MinPrice)
		maxPriceDiff := current.MinPrice * 3 / 10
		if maxPriceDiff < 2000 {
			maxPriceDiff = 2000
		}
		if priceDiff <= maxPriceDiff {

			score += 15 - (priceDiff * 5 / maxPriceDiff)
		}

		for _, tag := range cand.Tags {
			for _, curTag := range current.Tags {
				if tag.Ru == curTag.Ru || tag.En == curTag.En {
					score += 5
					break
				}
			}
		}

		if cand.Availability == "available" {
			score += 5
		}
		scored[i] = scoredItem{Item: cand, Score: score}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	result := make([]ItemFlatten, 0, limit)
	for i := 0; i < limit && i < len(scored); i++ {
		result = append(result, scored[i].Item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (c *Cache) MainPageNewHandler(w http.ResponseWriter, r *http.Request) {
	data := c.GetMainPageNew()
	if data == nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func parseFilterParams(r *http.Request) FilterParams {
	q := r.URL.Query()
	params := FilterParams{
		Page:   getInt(q, "page", 1),
		Limit:  getInt(q, "limit", 20),
		Sort:   q.Get("sort"),
		Search: q.Get("search"),
		Lang:   q.Get("lang"),
	}
	if params.Lang == "" {
		params.Lang = "ru"
	}

	if minStr := q.Get("priceMin"); minStr != "" {
		if min, err := strconv.Atoi(minStr); err == nil {
			params.PriceMin = &min
		}
	}
	if maxStr := q.Get("priceMax"); maxStr != "" {
		if max, err := strconv.Atoi(maxStr); err == nil {
			params.PriceMax = &max
		}
	}

	params.Brands = q["brands[]"]
	params.Countries = q["countries[]"]
	params.Materials = q["materials[]"]
	params.Categories = q["categories[]"]
	params.Types = q["types[]"]
	params.Colors = q["colors[]"]
	params.Tags = q["tags[]"]
	params.Sizes = q["sizes[]"]
	params.Availability = q["availability[]"]

	return params
}

func getInt(q map[string][]string, key string, def int) int {
	if vals, ok := q[key]; ok && len(vals) > 0 {
		if val, err := strconv.Atoi(vals[0]); err == nil {
			return val
		}
	}
	return def
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (c *Cache) CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post("https://3b7b2b24dfd8c527.mokky.dev/items", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		go c.Refresh()
	}
}

func (c *Cache) UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPatch, "https://3b7b2b24dfd8c527.mokky.dev/items/"+id, bytes.NewReader(bodyBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		go c.Refresh()
	}
}

func (c *Cache) DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, "https://3b7b2b24dfd8c527.mokky.dev/items/"+id, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		go c.Refresh()
	}
}

func (c *Cache) AdminItemsHandler(w http.ResponseWriter, r *http.Request) {
	items := c.GetItems()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (c *Cache) AdminItemHandler(w http.ResponseWriter, r *http.Request) {
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

func (c *Cache) ExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	exchangeRateMu.RLock()
	rate := exchangeRate
	exchangeRateMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"rate": rate})
}

func (c *Cache) FiltersHandler(w http.ResponseWriter, r *http.Request) {
	items := c.GetFlatItems()

	result := struct {
		Brands       []string `json:"brands"`
		Countries    []string `json:"countries"`
		Materials    []Lang   `json:"materials"`
		Categories   []Lang   `json:"categories"`
		Types        []Lang   `json:"types"`
		Colors       []Lang   `json:"colors"`
		Tags         []Lang   `json:"tags"`
		Sizes        []string `json:"sizes"`
		Availability []string `json:"availability"`
		MinPrice     int      `json:"minPrice"`
		MaxPrice     int      `json:"maxPrice"`
	}{}

	brandSet := make(map[string]bool)
	countrySet := make(map[string]bool)
	materialMap := make(map[string]Lang)
	categoryMap := make(map[string]Lang)
	typeMap := make(map[string]Lang)
	colorMap := make(map[string]Lang)
	tagMap := make(map[string]Lang)
	sizeSet := make(map[string]bool)
	availSet := make(map[string]bool)

	for _, item := range items {
		if item.Brand != "" {
			brandSet[item.Brand] = true
		}
		if item.Country.Ru != "" {
			countrySet[item.Country.Ru] = true
		}

		keyCat := item.Category.Ru + "|" + item.Category.En
		if _, ok := categoryMap[keyCat]; !ok && (item.Category.Ru != "" || item.Category.En != "") {
			categoryMap[keyCat] = item.Category
		}

		keyType := item.Type.Ru + "|" + item.Type.En
		if _, ok := typeMap[keyType]; !ok && (item.Type.Ru != "" || item.Type.En != "") {
			typeMap[keyType] = item.Type
		}

		for _, mat := range item.Structure {
			keyMat := mat.Name.Ru + "|" + mat.Name.En
			if _, ok := materialMap[keyMat]; !ok && (mat.Name.Ru != "" || mat.Name.En != "") {
				materialMap[keyMat] = mat.Name
			}
		}

		for _, col := range item.Color {
			keyCol := col.Ru + "|" + col.En
			if _, ok := colorMap[keyCol]; !ok && (col.Ru != "" || col.En != "") {
				colorMap[keyCol] = col
			}
		}

		for _, tag := range item.Tags {
			keyTag := tag.Ru + "|" + tag.En
			if _, ok := tagMap[keyTag]; !ok && (tag.Ru != "" || tag.En != "") {
				tagMap[keyTag] = tag
			}
		}

		for _, sz := range item.Sizes {
			sizeStr := toString(sz.Size)
			if sizeStr != "" {
				sizeSet[sizeStr] = true
			}
		}

		if item.Availability != "" {
			availSet[item.Availability] = true
		}
	}

	for b := range brandSet {
		result.Brands = append(result.Brands, b)
	}
	for c := range countrySet {
		result.Countries = append(result.Countries, c)
	}
	for _, v := range materialMap {
		result.Materials = append(result.Materials, v)
	}
	for _, v := range categoryMap {
		result.Categories = append(result.Categories, v)
	}
	for _, v := range typeMap {
		result.Types = append(result.Types, v)
	}
	for _, v := range colorMap {
		result.Colors = append(result.Colors, v)
	}
	for _, v := range tagMap {
		result.Tags = append(result.Tags, v)
	}
	for s := range sizeSet {
		result.Sizes = append(result.Sizes, s)
	}
	for a := range availSet {
		result.Availability = append(result.Availability, a)
	}

	minPrice := int(^uint(0) >> 1)
	maxPrice := 0
	for _, item := range items {
		if item.MinPrice < minPrice {
			minPrice = item.MinPrice
		}
		if item.MinPrice > maxPrice {
			maxPrice = item.MinPrice
		}
	}
	if minPrice == int(^uint(0)>>1) {
		minPrice = 0
	}
	result.MinPrice = minPrice
	result.MaxPrice = maxPrice

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
