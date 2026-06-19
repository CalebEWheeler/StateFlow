package shared

type OrderRequestBody struct {
	Address    Address `json:"address"`
	Currency   string  `json:"currency"`
	CustomerID string  `json:"customer_id"`
	Email      string  `json:"email"`
	Items      []Item  `json:"items"`
}

type Address struct {
	AdministrativeArea string `json:"administrative_area"`
	City               string `json:"city"`
	Country            string `json:"country"`
	Street             string `json:"street"`
}

type Item struct {
	ID       string  `json:"id"`
	Quantity int     `json:"quantity"`
	SKU      string  `json:"sku"`
	MSRP     float64 `json:"msrp"`
	Price    float64 `json:"price"`
}
