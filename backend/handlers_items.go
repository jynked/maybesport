package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type FilterParams struct {
	Page         int
	Limit        int
	Sort         string
	Search       string
	PriceMin     *int64
	PriceMax     *int64
	Brands       []string
	Countries    []string
	Materials    []string
	Categories   []string
	Types        []string
	Colors       []string
	Tags         []string
	Sizes        []string
	Availability []string
	Lang         string
}

func (c *DBCache) ItemsHandler(w http.ResponseWriter, r *http.Request) {
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

func (c *DBCache) ItemHandler(w http.ResponseWriter, r *http.Request) {
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

func (c *DBCache) SimilarHandler(w http.ResponseWriter, r *http.Request) {
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
		priceDiff := abs(current.MinPrice - cand.MinPrice)
		maxPriceDiff := current.MinPrice * 3 / 10
		if maxPriceDiff < 2000 {
			maxPriceDiff = 2000
		}
		if priceDiff <= maxPriceDiff {
			score += int(15 - (priceDiff * 5 / maxPriceDiff))
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

	for i := 0; i < len(scored)-1; i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].Score > scored[i].Score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	result := make([]ItemFlatten, 0, limit)
	for i := 0; i < limit && i < len(scored); i++ {
		result = append(result, scored[i].Item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (c *DBCache) MainPageNewHandler(w http.ResponseWriter, r *http.Request) {
	data := c.GetMainPageNew()
	if data == nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (c *DBCache) FiltersHandler(w http.ResponseWriter, r *http.Request) {
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
		MinPrice     int64    `json:"minPrice"`
		MaxPrice     int64    `json:"maxPrice"`
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

	var minPrice int64 = 1 << 62
	var maxPrice int64
	for _, item := range items {
		if item.MinPrice < minPrice {
			minPrice = item.MinPrice
		}
		if item.MinPrice > maxPrice {
			maxPrice = item.MinPrice
		}
	}
	if minPrice == 1<<62 {
		minPrice = 0
	}
	result.MinPrice = minPrice
	result.MaxPrice = maxPrice

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
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
		if min, err := strconv.ParseInt(minStr, 10, 64); err == nil {
			params.PriceMin = &min
		}
	}
	if maxStr := q.Get("priceMax"); maxStr != "" {
		if max, err := strconv.ParseInt(maxStr, 10, 64); err == nil {
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

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func buildAppliedFilters(params FilterParams) []AppliedFilter {
	var filters []AppliedFilter
	if params.PriceMin != nil || params.PriceMax != nil {
		minVal := "0"
		if params.PriceMin != nil {
			minVal = strconv.FormatInt(*params.PriceMin, 10)
		}
		maxVal := "∞"
		if params.PriceMax != nil {
			maxVal = strconv.FormatInt(*params.PriceMax, 10)
		}
		filters = append(filters, AppliedFilter{
			Key:   "price",
			Value: fmt.Sprintf("%s - %s", minVal, maxVal),
			Label: "Цена",
		})
	}
	addFilter := func(key, label string, values []string) {
		if len(values) > 0 {
			filters = append(filters, AppliedFilter{Key: key, Value: values, Label: label})
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

func filterItems(items []ItemFlatten, params FilterParams) []ItemFlatten {
	filtered := make([]ItemFlatten, 0, len(items))

	for _, item := range items {
		if !matchesPrice(item, params.PriceMin, params.PriceMax) {
			continue
		}
		if !matchesBrand(item, params.Brands) {
			continue
		}
		if !matchesCountry(item, params.Countries) {
			continue
		}
		if !matchesMaterials(item, params.Materials) {
			continue
		}
		if !matchesCategory(item, params.Categories) {
			continue
		}
		if !matchesType(item, params.Types) {
			continue
		}
		if !matchesColors(item, params.Colors) {
			continue
		}
		if !matchesTags(item, params.Tags) {
			continue
		}
		if !matchesSizes(item, params.Sizes) {
			continue
		}
		if !matchesAvailability(item, params.Availability) {
			continue
		}
		if !matchesSearch(item, params.Search, params.Lang) {
			continue
		}
		filtered = append(filtered, item)
	}

	sortFilteredItems(filtered, params.Sort, params.Lang)
	return filtered
}

func matchesPrice(item ItemFlatten, min, max *int64) bool {
	if min != nil && item.MinPrice < *min {
		return false
	}
	if max != nil && item.MinPrice > *max {
		return false
	}
	return true
}

func matchesBrand(item ItemFlatten, brands []string) bool {
	if len(brands) == 0 {
		return true
	}
	for _, b := range brands {
		if item.Brand == b {
			return true
		}
	}
	return false
}

func matchesCountry(item ItemFlatten, countries []string) bool {
	if len(countries) == 0 {
		return true
	}
	for _, c := range countries {
		if item.Country.Ru == c || item.Country.En == c {
			return true
		}
	}
	return false
}

func matchesMaterials(item ItemFlatten, materials []string) bool {
	if len(materials) == 0 {
		return true
	}
	for _, m := range materials {
		for _, mat := range item.Structure {
			if mat.Name.Ru == m || mat.Name.En == m {
				return true
			}
		}
	}
	return false
}

func matchesCategory(item ItemFlatten, categories []string) bool {
	if len(categories) == 0 {
		return true
	}
	for _, c := range categories {
		if item.Category.Ru == c || item.Category.En == c {
			return true
		}
	}
	return false
}

func matchesType(item ItemFlatten, types []string) bool {
	if len(types) == 0 {
		return true
	}
	for _, t := range types {
		if item.Type.Ru == t || item.Type.En == t {
			return true
		}
	}
	return false
}

func matchesColors(item ItemFlatten, colors []string) bool {
	if len(colors) == 0 {
		return true
	}
	for _, col := range colors {
		for _, c := range item.Color {
			if c.Ru == col || c.En == col {
				return true
			}
		}
	}
	return false
}

func matchesTags(item ItemFlatten, tags []string) bool {
	if len(tags) == 0 {
		return true
	}
	for _, tag := range tags {
		for _, t := range item.Tags {
			if t.Ru == tag || t.En == tag {
				return true
			}
		}
	}
	return false
}

func matchesSizes(item ItemFlatten, sizes []string) bool {
	if len(sizes) == 0 {
		return true
	}
	for _, sz := range sizes {
		for _, s := range item.Sizes {
			if toString(s.Size) == sz {
				return true
			}
		}
	}
	return false
}

func matchesAvailability(item ItemFlatten, availability []string) bool {
	if len(availability) == 0 {
		return true
	}
	for _, a := range availability {
		if item.Availability == a {
			return true
		}
	}
	return false
}

func matchesSearch(item ItemFlatten, search, lang string) bool {
	if search == "" {
		return true
	}
	search = strings.ToLower(search)
	title := ""
	if lang == "ru" {
		title = strings.ToLower(item.Title.Ru)
	} else {
		title = strings.ToLower(item.Title.En)
	}
	if strings.Contains(title, search) {
		return true
	}
	for _, tag := range item.Tags {
		tagText := ""
		if lang == "ru" {
			tagText = strings.ToLower(tag.Ru)
		} else {
			tagText = strings.ToLower(tag.En)
		}
		if strings.Contains(tagText, search) {
			return true
		}
	}
	return false
}

func sortFilteredItems(items []ItemFlatten, sortBy, lang string) {
	switch sortBy {
	case "price-asc":
		sort.Slice(items, func(i, j int) bool {
			return items[i].MinPrice < items[j].MinPrice
		})
	case "price-desc":
		sort.Slice(items, func(i, j int) bool {
			return items[i].MinPrice > items[j].MinPrice
		})
	case "alphabet-asc":
		sort.Slice(items, func(i, j int) bool {
			return getTitleByLang(items[i].Title, lang) < getTitleByLang(items[j].Title, lang)
		})
	case "alphabet-desc":
		sort.Slice(items, func(i, j int) bool {
			return getTitleByLang(items[i].Title, lang) > getTitleByLang(items[j].Title, lang)
		})
	}
}

func getTitleByLang(title Lang, lang string) string {
	if lang == "ru" {
		return title.Ru
	}
	return title.En
}
