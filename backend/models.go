package main

import "time"

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

type FavouriteItem struct {
	UniqueId string      `json:"uniqueId"`
	Size     interface{} `json:"size"`
}

type FavouriteItemResponse struct {
	UniqueId     string      `json:"uniqueId"`
	ID           int         `json:"id"`
	Title        Lang        `json:"title"`
	Image        string      `json:"image"`
	Size         interface{} `json:"size"`
	Price        int         `json:"price"`
	IsOnRequest  bool        `json:"isOnRequest"`
	Quantity     int         `json:"quantity"`
	Availability string      `json:"availability"`
}

type OrderStatus string

const (
	OrderStatusCreated    OrderStatus = "created"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusReceived   OrderStatus = "received"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type OrderStatusHistoryEntry struct {
	Status      OrderStatus `json:"status"`
	Timestamp   time.Time   `json:"timestamp"`
	Description string      `json:"description"`
}

type OrderItem struct {
	UniqueId string      `json:"uniqueId"`
	Size     interface{} `json:"size"`
	Price    int         `json:"price"`
	Quantity int         `json:"quantity"`
}

type Order struct {
	ID              int                       `json:"id"`
	CreatedAt       time.Time                 `json:"createdAt"`
	UpdatedAt       time.Time                 `json:"updatedAt"`
	Status          OrderStatus               `json:"status"`
	StatusHistory   []OrderStatusHistoryEntry `json:"statusHistory"`
	Items           []OrderItem               `json:"items"`
	TotalAmount     int                       `json:"totalAmount"`
	DeliveryAddress string                    `json:"deliveryAddress,omitempty"`
}

type CartItem struct {
	UniqueId string      `json:"uniqueId"`
	Size     interface{} `json:"size"`
	Quantity int         `json:"quantity"`
}

type User struct {
	ID         int             `json:"id"`
	Email      string          `json:"email"`
	Password   string          `json:"password,omitempty"`
	Favourites []FavouriteItem `json:"favourites"`
	Cart       []CartItem      `json:"cart"`
	Orders     []Order         `json:"orders"`
}
