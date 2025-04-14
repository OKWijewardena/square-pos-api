package square

import (
	"os"
	"net/http"
	"net/http/httptest"
	"testing"
)

// This test simulates the Square server using httptest.Server
func TestCreateSquareOrder(t *testing.T) {
	// Step 1: Start a fake HTTP server to simulate Square API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"order": {
				"id": "mock-order-id",
				"reference_id": "10",
				"created_at": "2025-04-14T10:00:00Z",
				"state": "OPEN",
				"line_items": [
					{
						"name": "Burger",
						"quantity": "2",
						"base_price_money": { "amount": 1200 },
						"total_money": { "amount": 2400 }
					}
				],
				"total_money": { "amount": 2400 },
				"total_tax_money": { "amount": 0 },
				"total_discount_money": { "amount": 0 }
			}
		}`))
	}))
	defer mockServer.Close()

	// Replace in test by modifying the env or faking the host
	os.Setenv("SQUARE_ACCESS_TOKEN", "EAAAl7pUMdR4fYWjgaXg1xvdy-jhH1AEq0uN-WF6zDyRPF7FeCum6OKkaIwMSwcj")
	os.Setenv("SQUARE_LOCATION_ID", "LNE1FE7JFW9PJ")


	// Direct test of CallSquareAPI using mockServer
	body := map[string]interface{}{
		"order": map[string]interface{}{
			"location_id":  "LNE1FE7JFW9PJ",
			"reference_id": "10",
			"line_items": []map[string]interface{}{
				{
					"name":     "Burger",
					"quantity": "2",
					"base_price_money": map[string]interface{}{
						"amount":   1200,
						"currency": "USD",
					},
				},
			},
		},
		"idempotency_key": "test-key-123",
	}

	resp, err := CallSquareAPI("POST", mockServer.URL, body)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(resp) == 0 {
		t.Error("Expected response from mock Square server, got empty")
	}
}
