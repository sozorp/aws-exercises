package main

type Order struct {
	OrderID string `json:"order_id"`
}

type Output struct {
	Status  string `json:"status"`
	OrderID string `json:"order_id"`
}
