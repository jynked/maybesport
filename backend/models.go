package main

type Item struct {
	ID          int    `json:"id"`
	Type        Lang   `json:"type"`
	Title       Lang   `json:"title"`
	Description Lang   `json:"description"`
	Brand       string `json:"brand"`
	Country     Lang   `json:"country"`
	Structure   []struct {
		Name    Lang `json:"name"`
		Percent int  `json:"percent"`
	} `json:"structure"`
	Category  Lang      `json:"category"`
	CreatedAt string    `json:"createdAt"`
	Items     []SubItem `json:"items"`
}

type SubItem struct {
	UniqueId string   `json:"uniqueId"`
	Images   []string `json:"images"`
	Color    []Lang   `json:"color"`
	Tags     []Lang   `json:"tags"`
	Sizes    []Size   `json:"sizes"`
}

type Size struct {
	Size        interface{} `json:"size"`
	Price       int         `json:"price"`
	IsOnRequest bool        `json:"isOnRequest"`
	Quantity    int         `json:"quantity"`
}

type Lang struct {
	Ru string `json:"ru"`
	En string `json:"en"`
}

type ItemFlatten struct {
	UniqueId    string `json:"uniqueId"`
	ID          int    `json:"id"`
	Type        Lang   `json:"type"`
	Title       Lang   `json:"title"`
	Description Lang   `json:"description"`
	Brand       string `json:"brand"`
	Country     Lang   `json:"country"`
	Structure   []struct {
		Name    Lang `json:"name"`
		Percent int  `json:"percent"`
	} `json:"structure"`
	Category  Lang     `json:"category"`
	CreatedAt string   `json:"createdAt"`
	Images    []string `json:"images"`
	Color     []Lang   `json:"color"`
	Tags      []Lang   `json:"tags"`
	Sizes     []Size   `json:"sizes"`

	Availability  string        `json:"availability"`
	MinPrice      int           `json:"minPrice"`
	TotalQuantity int           `json:"totalQuantity"`
	SiblingItems  []SiblingItem `json:"siblingItems,omitempty"`
}

type SiblingItem struct {
	UniqueId string `json:"uniqueId"`
	Image    string `json:"image"`
}

type ItemsResponse struct {
	Items          []ItemFlatten   `json:"items"`
	Total          int             `json:"total"`
	Page           int             `json:"page"`
	Limit          int             `json:"limit"`
	TotalPages     int             `json:"totalPages"`
	AppliedFilters []AppliedFilter `json:"appliedFilters"`
}

type AppliedFilter struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Label string      `json:"label"`
}

type User struct {
    ID        int      `json:"id"`
    Email     string   `json:"email"`
    Password  string   `json:"password,omitempty"`
    Favourites []string `json:"favourites"`
}