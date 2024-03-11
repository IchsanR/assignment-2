package core

import "time"

type Item struct {
	ItemCode    string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

type Order struct {
	OrderedAt    time.Time `json:"orderedAt"`
	CustomerName string    `json:"customerName"`
	Items        []Item    `json:"items"`
}

type Products interface {
	InsertOrder(order Order) error
}
