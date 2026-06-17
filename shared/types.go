package shared

type OrderRequestBody struct {
	CustomerID string  `json:"customer_id"`
	Email      string  `json:"email"`
	Address    Address `json:"address"`
	Items      []Item  `json:"items"`
	Currency   string  `json:"currency"`
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
}

type Item struct {
	ID       string `json:"id"`
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}
