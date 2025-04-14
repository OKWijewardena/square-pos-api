package square

import (
	"os"

	"github.com/google/uuid"
)

type OrderItem struct {
	Name      string
	Quantity  string
	BasePrice int
}

func CreateSquareOrder(tableNumber string, items []OrderItem) ([]byte, error) {
	url := "https://connect.squareupsandbox.com/v2/orders"

	// Get location ID from env at runtime
	squareLocationID := os.Getenv("SQUARE_LOCATION_ID")

	// Make the idempotency key unique for every request
	idempotencyKey := uuid.New().String()

	// Build the Square order body
	body := map[string]interface{}{
		"order": map[string]interface{}{
			"location_id":  squareLocationID,
			"reference_id": tableNumber,
			"line_items":   []map[string]interface{}{},
		},
		"idempotency_key": idempotencyKey,
	}

	// Add each item to the order
	for _, item := range items {
		line := map[string]interface{}{
			"name":     item.Name,
			"quantity": item.Quantity,
			"base_price_money": map[string]interface{}{
				"amount":   item.BasePrice, // in cents
				"currency": "USD",
			},
		}
		body["order"].(map[string]interface{})["line_items"] =
			append(body["order"].(map[string]interface{})["line_items"].([]map[string]interface{}), line)
	}

	return CallSquareAPI("POST", url, body)
}
