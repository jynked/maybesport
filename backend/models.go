package main

import "time"

type Lang struct {
	Ru string `json:"ru"`
	En string `json:"en"`
}

type Size struct {
	Size        interface{} `json:"size"`
	Price       int64       `json:"price"`
	IsOnRequest bool        `json:"isOnRequest"`
	Quantity    int64       `json:"quantity"`
	PriceCny    int64       `json:"priceCny"`
}

type SubItem struct {
	UniqueId string   `json:"uniqueId"`
	Images   []string `json:"images"`
	Color    []Lang   `json:"color"`
	Tags     []Lang   `json:"tags"`
	Sizes    []Size   `json:"sizes"`
}

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
	CreatedAt time.Time `json:"createdAt"`
	Items     []SubItem `json:"items"`
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
	Category      Lang          `json:"category"`
	CreatedAt     time.Time     `json:"createdAt"`
	Images        []string      `json:"images"`
	Color         []Lang        `json:"color"`
	Tags          []Lang        `json:"tags"`
	Sizes         []Size        `json:"sizes"`
	Availability  string        `json:"availability"`
	MinPrice      int64         `json:"minPrice"`
	TotalQuantity int64         `json:"totalQuantity"`
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
	Price        int64       `json:"price"`
	IsOnRequest  bool        `json:"isOnRequest"`
	Quantity     int64       `json:"quantity"`
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
	ID       int64       `json:"id"`
	UniqueId string      `json:"uniqueId"`
	Size     interface{} `json:"size"`
	Price    int64       `json:"price"`
	Quantity int64       `json:"quantity"`
	Status   OrderStatus `json:"status"`
}

type OrderItemStatusHistoryEntry struct {
	Status      OrderStatus `json:"status"`
	Timestamp   time.Time   `json:"timestamp"`
	Description string      `json:"description"`
}

type OrderItemWithHistory struct {
	OrderItem
	StatusHistory []OrderItemStatusHistoryEntry `json:"statusHistory"`
}

type Order struct {
	ID              int                       `json:"id"`
	CreatedAt       time.Time                 `json:"createdAt"`
	UpdatedAt       time.Time                 `json:"updatedAt"`
	Status          OrderStatus               `json:"status"`
	StatusHistory   []OrderStatusHistoryEntry `json:"statusHistory"`
	Items           []OrderItem               `json:"items"`
	TotalAmount     int64                     `json:"totalAmount"`
	DeliveryAddress *string                   `json:"deliveryAddress,omitempty"`
}

type CartItemResponse struct {
	UniqueId     string      `json:"uniqueId"`
	ID           int         `json:"id"`
	Title        Lang        `json:"title"`
	Image        string      `json:"image"`
	Size         interface{} `json:"size"`
	Price        int64       `json:"price"`
	Quantity     int64       `json:"quantity"`
	IsOnRequest  bool        `json:"isOnRequest"`
	Stock        int64       `json:"stock"`
	Color        Lang        `json:"color"`
	Availability string      `json:"availability"`
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type AdminOrderListItem struct {
	ID              int         `json:"id"`
	UserID          int         `json:"user_id"`
	UserEmail       string      `json:"user_email"`
	Status          OrderStatus `json:"status"`
	TotalAmount     int64       `json:"total_amount"`
	DeliveryAddress *string     `json:"delivery_address"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}
