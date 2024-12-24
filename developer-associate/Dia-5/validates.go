package main

import (
	"encoding/json"
	"fmt"
)

func ValidateOrder(input json.RawMessage) (*Order, error) {
	var order Order

	if err := json.Unmarshal(input, &order); err != nil {
		return &Order{}, fmt.Errorf("The input is not valid for the flow")
	}

	return &order, nil
}
