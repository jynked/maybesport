package main

import (
    "sort"
    "strconv"
    "strings"
)

type FilterParams struct {
    Page         int
    Limit        int
    Sort         string
    Search       string
    PriceMin     *int
    PriceMax     *int
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

func matchesPrice(item ItemFlatten, min, max *int) bool {
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

func toString(v interface{}) string {
    switch val := v.(type) {
    case string:
        return val
    case int:
        return strconv.Itoa(val)
    case float64:
        return strconv.FormatFloat(val, 'f', -1, 64)
    default:
        return ""
    }
}