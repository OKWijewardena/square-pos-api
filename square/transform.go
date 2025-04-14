package square

import (
	"fmt"
)

type CustomOrderResponse struct {
	ID       string             `json:"id"`
	OpenedAt string             `json:"opened_at"`
	IsClosed bool               `json:"is_closed"`
	Table    string             `json:"table"`
	Items    []CustomOrderItem  `json:"items"`
	Totals   CustomOrderTotals  `json:"totals"`
}

type CustomOrderItem struct {
	Name      string               `json:"name"`
	Comment   string               `json:"comment"`
	UnitPrice int                  `json:"unit_price"`
	Quantity  int                  `json:"quantity"`
	Discounts []CustomItemDiscount `json:"discounts"`
	Modifiers []CustomItemModifier `json:"modifiers"`
	Amount    int                  `json:"amount"`
}

type CustomItemDiscount struct {
	Name         string `json:"name"`
	IsPercentage bool   `json:"is_percentage"`
	Value        int    `json:"value"`
	Amount       int    `json:"amount"`
}

type CustomItemModifier struct {
	Name      string `json:"name"`
	UnitPrice int    `json:"unit_price"`
	Quantity  int    `json:"quantity"`
	Amount    int    `json:"amount"`
}

type CustomOrderTotals struct {
	Discounts     int `json:"discounts"`
	Due           int `json:"due"`
	Tax           int `json:"tax"`
	ServiceCharge int `json:"service_charge"`
	Paid          int `json:"paid"`
	Tips          int `json:"tips"`
	Total         int `json:"total"`
}

func TransformSquareOrder(squareOrder map[string]interface{}) CustomOrderResponse {
	order := squareOrder["order"].(map[string]interface{})

	// Prepare item list
	var items []CustomOrderItem

	if rawItems, ok := order["line_items"].([]interface{}); ok {
		for _, raw := range rawItems {
			item := raw.(map[string]interface{})

			// Parse quantity string to int
			qty := 1
			if q, ok := item["quantity"].(string); ok {
				fmt.Sscanf(q, "%d", &qty)
			}

			price := 0
			if bp, ok := item["base_price_money"].(map[string]interface{}); ok {
				price = int(bp["amount"].(float64))
			}

			amount := 0
			if total, ok := item["total_money"].(map[string]interface{}); ok {
				amount = int(total["amount"].(float64))
			}

			// Build modifiers
			var modifiers []CustomItemModifier
			if rawMods, ok := item["modifiers"].([]interface{}); ok {
				for _, m := range rawMods {
					mod := m.(map[string]interface{})
					q := 1
					if qs, ok := mod["quantity"].(string); ok {
						fmt.Sscanf(qs, "%d", &q)
					}

					unitPrice := 0
					if modMoney, ok := mod["base_price_money"].(map[string]interface{}); ok {
						unitPrice = int(modMoney["amount"].(float64))
					}

					amount := unitPrice * q
					modifiers = append(modifiers, CustomItemModifier{
						Name:      mod["name"].(string),
						Quantity:  q,
						UnitPrice: unitPrice,
						Amount:    amount,
					})
				}
			}

			// Build discounts
			var discounts []CustomItemDiscount
			if rawDiscounts, ok := item["discounts"].([]interface{}); ok {
				for _, d := range rawDiscounts {
					disc := d.(map[string]interface{})
					amount := int(disc["amount_money"].(map[string]interface{})["amount"].(float64))
					value := 0
					if v, ok := disc["percentage"].(string); ok {
						fmt.Sscanf(v, "%d", &value)
					} else if v, ok := disc["amount_money"].(map[string]interface{}); ok {
						value = int(v["amount"].(float64))
					}

					discounts = append(discounts, CustomItemDiscount{
						Name:         disc["name"].(string),
						IsPercentage: disc["type"].(string) == "FIXED_PERCENTAGE",
						Value:        value,
						Amount:       amount,
					})
				}
			}

			customItem := CustomOrderItem{
				Name:      item["name"].(string),
				Comment:   "", // Square doesn't include comments in line_items
				UnitPrice: price,
				Quantity:  qty,
				Amount:    amount,
				Modifiers: modifiers,
				Discounts: discounts,
			}

			items = append(items, customItem)
		}
	}

	// Build totals
	discounts := 0
	if d, ok := order["total_discount_money"].(map[string]interface{}); ok {
		discounts = int(d["amount"].(float64))
	}

	tax := 0
	if t, ok := order["total_tax_money"].(map[string]interface{}); ok {
		tax = int(t["amount"].(float64))
	}

	total := 0
	if t, ok := order["total_money"].(map[string]interface{}); ok {
		total = int(t["amount"].(float64))
	}

	return CustomOrderResponse{
		ID:       order["id"].(string),
		OpenedAt: order["created_at"].(string),
		IsClosed: order["state"].(string) == "COMPLETED",
		Table:    order["reference_id"].(string),
		Items:    items,
		Totals: CustomOrderTotals{
			Discounts:     discounts,
			Tax:           tax,
			ServiceCharge: 0,
			Due:           0,
			Paid:          0,
			Tips:          0,
			Total:         total,
		},
	}
}
